package gateway

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// mockForwarder adalah Forwarder mock untuk test.
type mockForwarder struct {
	result ForwardResult
	err    error
	calls  int
}

func (m *mockForwarder) Forward(_ context.Context, _ string, _ []byte) (ForwardResult, error) {
	m.calls++
	return m.result, m.err
}

// mockLogStore adalah LogStore mock untuk test.
type mockLogStore struct {
	logs  []model.RequestLog
	err   error
	calls int
}

func (m *mockLogStore) Insert(_ context.Context, log model.RequestLog) (int64, error) {
	m.calls++
	if m.err != nil {
		return 0, m.err
	}
	m.logs = append(m.logs, log)
	return int64(m.calls), nil
}

var (
	adminClaims = auth.Claims{
		Username: "admin",
		Role:     model.RoleAdminGateway,
		AppName:  "API Gateway",
		RegisteredClaims: jwt.RegisteredClaims{Subject: "1"},
	}
	appUserClaims = auth.Claims{
		Username: "marketplace",
		Role:     model.RoleAppUser,
		AppName:  "Marketplace",
		RegisteredClaims: jwt.RegisteredClaims{Subject: "2"},
	}
	monitoringClaims = auth.Claims{
		Username: "insight",
		Role:     model.RoleMonitoringUser,
		AppName:  "UMKM Insight",
		RegisteredClaims: jwt.RegisteredClaims{Subject: "3"},
	}
)

func newTestService(fwd Forwarder, logs LogStore) *Service {
	return New(fwd, logs, UpstreamConfig{
		SmartBankURL:   "http://smartbank.local",
		MarketplaceURL: "http://marketplace.local",
		LogisticsURL:   "http://logistics.local",
		SupplierHubURL: "http://supplierhub.local",
	}, time.Now)
}

func TestService_ForwardPayment_Success(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{"status":"ok"}`), DurationMS: 50}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{FromApp: "Marketplace", FromUser: "user1", ToUser: "user2", Amount: 10000, ServiceType: "payment"}
	resp, err := svc.ForwardPayment(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Forwarded {
		t.Fatal("expected forwarded=true")
	}
	if resp.Upstream != "smartbank" {
		t.Fatalf("expected upstream=smartbank, got %s", resp.Upstream)
	}
	if fwd.calls != 1 {
		t.Fatalf("expected 1 forward call, got %d", fwd.calls)
	}
	if logs.calls != 1 {
		t.Fatalf("expected 1 log call, got %d", logs.calls)
	}
}

func TestService_ForwardPayment_AdminAccess(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{}`), DurationMS: 10}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{FromApp: "Test", FromUser: "admin", ToUser: "user", Amount: 100, ServiceType: "payment"}
	_, err := svc.ForwardPayment(context.Background(), adminClaims, req)
	if err != nil {
		t.Fatalf("admin should have access: %v", err)
	}
}

func TestService_ForwardPayment_MonitoringForbidden(t *testing.T) {
	fwd := &mockForwarder{}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{FromApp: "Test", FromUser: "user1", ToUser: "user2", Amount: 100, ServiceType: "payment"}
	_, err := svc.ForwardPayment(context.Background(), monitoringClaims, req)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if fwd.calls != 0 {
		t.Fatal("forwarder should not be called for forbidden role")
	}
}

func TestService_ForwardPayment_InvalidPayload(t *testing.T) {
	fwd := &mockForwarder{}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{} // all fields empty
	_, err := svc.ForwardPayment(context.Background(), appUserClaims, req)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
	if fwd.calls != 0 {
		t.Fatal("forwarder should not be called for invalid payload")
	}
}

func TestService_ForwardPayment_UpstreamNotConfigured(t *testing.T) {
	fwd := &mockForwarder{err: ErrUpstreamNotConfigured}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{FromApp: "Marketplace", FromUser: "user1", ToUser: "user2", Amount: 10000, ServiceType: "payment"}
	resp, err := svc.ForwardPayment(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("not configured should not error: %v", err)
	}
	if resp.Forwarded {
		t.Fatal("expected forwarded=false")
	}
	if resp.Upstream != "not_configured" {
		t.Fatalf("expected upstream=not_configured, got %s", resp.Upstream)
	}
	if logs.calls != 1 {
		t.Fatalf("expected 1 log call, got %d", logs.calls)
	}
}

func TestService_ForwardPayment_UpstreamFailure(t *testing.T) {
	fwd := &mockForwarder{err: errors.New("upstream service failure: connection refused")}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{FromApp: "Marketplace", FromUser: "user1", ToUser: "user2", Amount: 10000, ServiceType: "payment"}
	_, err := svc.ForwardPayment(context.Background(), appUserClaims, req)
	if err == nil {
		t.Fatal("expected error for upstream failure")
	}
	if logs.calls != 1 {
		t.Fatalf("expected 1 log call even on failure, got %d", logs.calls)
	}
}

func TestService_ForwardSmartBank_Success(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{"balance":1000}`), DurationMS: 30}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := SmartBankRequest{Action: "check_balance"}
	resp, err := svc.ForwardSmartBank(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Forwarded {
		t.Fatal("expected forwarded=true")
	}
}

func TestService_ForwardSmartBank_InvalidAction(t *testing.T) {
	fwd := &mockForwarder{}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := SmartBankRequest{Action: ""}
	_, err := svc.ForwardSmartBank(context.Background(), appUserClaims, req)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestService_ForwardMarketplace_Success(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{"order_id":"ORD-1"}`), DurationMS: 25}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := MarketplaceRequest{Action: "get_order"}
	resp, err := svc.ForwardMarketplace(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Forwarded {
		t.Fatal("expected forwarded=true")
	}
	if resp.Upstream != "marketplace" {
		t.Fatalf("expected upstream=marketplace, got %s", resp.Upstream)
	}
}

func TestService_ForwardLogistics_Success(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{"delivery_id":"DEL-1"}`), DurationMS: 40}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := LogisticsRequest{OrderID: "ORD-1", Address: "Jl. Merdeka", Distance: 10, ShippingType: "express"}
	resp, err := svc.ForwardLogistics(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Forwarded {
		t.Fatal("expected forwarded=true")
	}
	if resp.Upstream != "logistikit" {
		t.Fatalf("expected upstream=logistikit, got %s", resp.Upstream)
	}
}

func TestService_ForwardLogistics_InvalidPayload(t *testing.T) {
	fwd := &mockForwarder{}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := LogisticsRequest{OrderID: "", Address: "", Distance: 0}
	_, err := svc.ForwardLogistics(context.Background(), appUserClaims, req)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestService_ForwardSupplier_Success(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{"order_id":"SUP-1"}`), DurationMS: 35}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := SupplierRequest{SupplierID: "SUP-1", Material: "Beras", Qty: 50, TotalCost: 250000}
	resp, err := svc.ForwardSupplier(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Forwarded {
		t.Fatal("expected forwarded=true")
	}
	if resp.Upstream != "supplierhub" {
		t.Fatalf("expected upstream=supplierhub, got %s", resp.Upstream)
	}
}

func TestService_ForwardSupplier_InvalidPayload(t *testing.T) {
	fwd := &mockForwarder{}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := SupplierRequest{} // all empty
	_, err := svc.ForwardSupplier(context.Background(), appUserClaims, req)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestService_NilLogs(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{}`), DurationMS: 10}}
	svc := New(fwd, nil, UpstreamConfig{SmartBankURL: "http://smartbank.local"}, time.Now)

	req := PaymentRequest{FromApp: "Test", FromUser: "u1", ToUser: "u2", Amount: 100, ServiceType: "payment"}
	_, err := svc.ForwardPayment(context.Background(), appUserClaims, req)
	if err != nil {
		t.Fatalf("should not panic with nil logs: %v", err)
	}
}

func TestService_LogContent(t *testing.T) {
	fwd := &mockForwarder{result: ForwardResult{StatusCode: 200, Body: []byte(`{"ok":true}`), DurationMS: 42}}
	logs := &mockLogStore{}
	svc := newTestService(fwd, logs)

	req := PaymentRequest{FromApp: "Marketplace", FromUser: "user1", ToUser: "user2", Amount: 5000, ServiceType: "payment"}
	svc.ForwardPayment(context.Background(), appUserClaims, req)

	if len(logs.logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logs.logs))
	}
	logEntry := logs.logs[0]
	if logEntry.Endpoint != "/gateway/payment" {
		t.Fatalf("expected endpoint=/gateway/payment, got %s", logEntry.Endpoint)
	}
	if logEntry.Method != "POST" {
		t.Fatalf("expected method=POST, got %s", logEntry.Method)
	}
	if logEntry.SourceApp != "Marketplace" {
		t.Fatalf("expected source_app=Marketplace, got %s", logEntry.SourceApp)
	}
	if logEntry.Status != 200 {
		t.Fatalf("expected status=200, got %d", logEntry.Status)
	}
	if logEntry.DurationMS == nil || *logEntry.DurationMS != 42 {
		t.Fatalf("expected duration_ms=42, got %v", logEntry.DurationMS)
	}
}
