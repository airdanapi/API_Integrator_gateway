package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/gateway"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

type stubGatewaySvc struct {
	paymentResp    gateway.GatewayResponse
	smartbankResp  gateway.GatewayResponse
	marketResp     gateway.GatewayResponse
	logisticsResp  gateway.GatewayResponse
	supplierResp   gateway.GatewayResponse
	err            error
	seenClaims     auth.Claims
	seenPayment    gateway.PaymentRequest
	seenSmartBank  gateway.SmartBankRequest
	seenMarket     gateway.MarketplaceRequest
	seenLogistics  gateway.LogisticsRequest
	seenSupplier   gateway.SupplierRequest
}

func (s *stubGatewaySvc) ForwardPayment(_ context.Context, claims auth.Claims, req gateway.PaymentRequest) (gateway.GatewayResponse, error) {
	s.seenClaims = claims
	s.seenPayment = req
	return s.paymentResp, s.err
}

func (s *stubGatewaySvc) ForwardSmartBank(_ context.Context, claims auth.Claims, req gateway.SmartBankRequest) (gateway.GatewayResponse, error) {
	s.seenClaims = claims
	s.seenSmartBank = req
	return s.smartbankResp, s.err
}

func (s *stubGatewaySvc) ForwardMarketplace(_ context.Context, claims auth.Claims, req gateway.MarketplaceRequest) (gateway.GatewayResponse, error) {
	s.seenClaims = claims
	s.seenMarket = req
	return s.marketResp, s.err
}

func (s *stubGatewaySvc) ForwardLogistics(_ context.Context, claims auth.Claims, req gateway.LogisticsRequest) (gateway.GatewayResponse, error) {
	s.seenClaims = claims
	s.seenLogistics = req
	return s.logisticsResp, s.err
}

func (s *stubGatewaySvc) ForwardSupplier(_ context.Context, claims auth.Claims, req gateway.SupplierRequest) (gateway.GatewayResponse, error) {
	s.seenClaims = claims
	s.seenSupplier = req
	return s.supplierResp, s.err
}

func makeGatewayVerifier(role model.Role, username, appName string) *stubVerifier {
	return &stubVerifier{claims: auth.Claims{Username: username, Role: role, AppName: appName}}
}

func TestGatewayRequiresAuthentication(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		GatewayService: &stubGatewaySvc{},
	})

	endpoints := []string{
		"/gateway/payment",
		"/gateway/smartbank",
		"/gateway/marketplace",
		"/gateway/logistics",
		"/gateway/supplier",
	}

	for _, ep := range endpoints {
		resp, err := app.Test(httptest.NewRequest(http.MethodPost, ep, bytes.NewBufferString(`{}`)))
		if err != nil {
			t.Fatalf("POST %s error: %v", ep, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("POST %s status = %d, want 401", ep, resp.StatusCode)
		}
	}
}

func TestGatewayPaymentInvalidJSON(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		GatewayService: &stubGatewaySvc{},
	})

	req := httptest.NewRequest(http.MethodPost, "/gateway/payment", bytes.NewBufferString(`not-json`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/payment error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestGatewayPaymentForbiddenMonitoring(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeGatewayVerifier(model.RoleMonitoringUser, "insight", "UMKM Insight"),
		GatewayService: &stubGatewaySvc{err: gateway.ErrForbidden},
	})

	body := `{"from_app":"Marketplace","from_user":"u1","to_user":"u2","amount":100,"service_type":"payment"}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/payment", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/payment error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", resp.StatusCode)
	}
}

func TestGatewayPaymentSuccess(t *testing.T) {
	svc := &stubGatewaySvc{
		paymentResp: gateway.GatewayResponse{
			Status:        "success",
			TransactionID: "gw-payment-123",
			Message:       "request forwarded successfully",
			Forwarded:     true,
			Upstream:      "smartbank",
		},
	}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		GatewayService: svc,
	})

	body := `{"from_app":"Marketplace","from_user":"user1","to_user":"user2","amount":10000,"service_type":"payment"}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/payment", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/payment error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result struct {
		Status string                  `json:"status"`
		Data   gateway.GatewayResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.Status != "success" {
		t.Fatalf("status = %q, want success", result.Status)
	}
	if !result.Data.Forwarded {
		t.Fatal("expected forwarded=true")
	}
	if svc.seenPayment.FromApp != "Marketplace" || svc.seenPayment.Amount != 10000 {
		t.Fatalf("seenPayment = %#v", svc.seenPayment)
	}
	if svc.seenClaims.Username != "marketplace" {
		t.Fatalf("seenClaims = %#v", svc.seenClaims)
	}
}

func TestGatewayPaymentInvalidPayload(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		GatewayService: &stubGatewaySvc{err: gateway.ErrInvalidPayload},
	})

	body := `{"from_app":"","from_user":"","to_user":"","amount":0,"service_type":""}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/payment", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/payment error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestGatewayPaymentUpstreamFailure(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		GatewayService: &stubGatewaySvc{err: gateway.ErrUpstreamFailure},
	})

	body := `{"from_app":"Marketplace","from_user":"u1","to_user":"u2","amount":100,"service_type":"payment"}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/payment", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/payment error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
}

func TestGatewaySmartBankSuccess(t *testing.T) {
	svc := &stubGatewaySvc{
		smartbankResp: gateway.GatewayResponse{
			Status:    "success",
			Forwarded: true,
			Upstream:  "smartbank",
		},
	}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "smartbank", "SmartBank"),
		GatewayService: svc,
	})

	body := `{"action":"check_balance","payload":{"account":"123"}}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/smartbank", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/smartbank error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if svc.seenSmartBank.Action != "check_balance" {
		t.Fatalf("seenSmartBank = %#v", svc.seenSmartBank)
	}
}

func TestGatewayMarketplaceSuccess(t *testing.T) {
	svc := &stubGatewaySvc{
		marketResp: gateway.GatewayResponse{
			Status:    "success",
			Forwarded: true,
			Upstream:  "marketplace",
		},
	}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAdminGateway, "admin", "API Gateway"),
		GatewayService: svc,
	})

	body := `{"action":"get_order"}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/marketplace", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/marketplace error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

func TestGatewayLogisticsSuccess(t *testing.T) {
	svc := &stubGatewaySvc{
		logisticsResp: gateway.GatewayResponse{
			Status:    "success",
			DeliveryID: "DEL-001",
			Forwarded: true,
			Upstream:  "logistikit",
		},
	}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
		GatewayService: svc,
	})

	body := `{"order_id":"ORD-1","address":"Jl. Merdeka","distance":10,"shipping_type":"express"}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/logistics", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/logistics error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if svc.seenLogistics.OrderID != "ORD-1" {
		t.Fatalf("seenLogistics = %#v", svc.seenLogistics)
	}
}

func TestGatewaySupplierSuccess(t *testing.T) {
	svc := &stubGatewaySvc{
		supplierResp: gateway.GatewayResponse{
			Status:    "success",
			OrderID:   "SUP-001",
			Forwarded: true,
			Upstream:  "supplierhub",
		},
	}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:  makeGatewayVerifier(model.RoleAppUser, "supplier", "SupplierHub"),
		GatewayService: svc,
	})

	body := `{"supplier_id":"SUP-1","material":"Beras","qty":50,"total_cost":250000}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/supplier", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/supplier error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if svc.seenSupplier.SupplierID != "SUP-1" || svc.seenSupplier.Qty != 50 {
		t.Fatalf("seenSupplier = %#v", svc.seenSupplier)
	}
}

func TestGatewayNilServiceReturnsInternalError(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier: makeGatewayVerifier(model.RoleAppUser, "marketplace", "Marketplace"),
	})

	body := `{"from_app":"Marketplace","from_user":"u1","to_user":"u2","amount":100,"service_type":"payment"}`
	req := httptest.NewRequest(http.MethodPost, "/gateway/payment", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /gateway/payment error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", resp.StatusCode)
	}
}
