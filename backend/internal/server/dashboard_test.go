package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/dashboard"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

// ─── stub DashboardService ──────────────────────────────────────────────────

type stubDashboardSvc struct {
	summary    dashboard.TrafficSummary
	indicators []dashboard.ServiceIndicator
	logs       []dashboard.AuditLogEntry
	totalLogs  int64
}

func (s *stubDashboardSvc) GetTrafficSummary(_ context.Context, _ time.Time) (dashboard.TrafficSummary, error) {
	return s.summary, nil
}
func (s *stubDashboardSvc) GetServiceIndicators(_ context.Context) ([]dashboard.ServiceIndicator, error) {
	return s.indicators, nil
}
func (s *stubDashboardSvc) GetAuditLogs(_ context.Context, _, _ int) ([]dashboard.AuditLogEntry, int64, error) {
	return s.logs, s.totalLogs, nil
}

// ─── stub TokenVerifier ─────────────────────────────────────────────────────

type stubVerifier struct {
	claims auth.Claims
}

func (sv *stubVerifier) Validate(_ string) (auth.Claims, error) {
	return sv.claims, nil
}

func makeVerifier(role model.Role, appName string) *stubVerifier {
	return &stubVerifier{claims: auth.Claims{
		Username: "testuser",
		Role:     role,
		AppName:  appName,
	}}
}

// ─── Tests ──────────────────────────────────────────────────────────────────

func TestDashboardAdminRequiresAuthentication(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		DashboardService: &stubDashboardSvc{},
		// TokenVerifier nil → requireToken akan menolak request tanpa token
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard/admin", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /dashboard/admin error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d (Unauthorized)", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestDashboardAdminForbidsNonAdminRole(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		DashboardService: &stubDashboardSvc{},
		TokenVerifier:    makeVerifier(model.RoleAppUser, "Marketplace"),
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard/admin", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /dashboard/admin error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want %d (Forbidden)", resp.StatusCode, http.StatusForbidden)
	}
}

func TestDashboardAdminForbidsMonitoringRole(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		DashboardService: &stubDashboardSvc{},
		TokenVerifier:    makeVerifier(model.RoleMonitoringUser, "UMKM Insight"),
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard/admin", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /dashboard/admin error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("monitoring role status = %d, want %d (Forbidden)", resp.StatusCode, http.StatusForbidden)
	}
}

func TestDashboardAdminReturnsContract(t *testing.T) {
	svc := &stubDashboardSvc{
		summary: dashboard.TrafficSummary{
			TotalRequests:  100,
			SuccessCount:   84,
			ErrorCount:     16,
			SuccessRatePct: 84.0,
			AvgDurationMS:  120,
		},
		indicators: []dashboard.ServiceIndicator{
			{AppName: "Marketplace", Status: "inactive", LastRequest: time.Now().Add(-8 * 24 * time.Hour)},
			{AppName: "POS", Status: "active", LastRequest: time.Now().Add(-1 * time.Hour)},
		},
		logs: []dashboard.AuditLogEntry{
			{ID: 1, SourceApp: "POS", Endpoint: "/gateway/payment", Method: "POST", Status: 200, Timestamp: time.Now()},
		},
		totalLogs: 50,
	}

	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		DashboardService: svc,
		TokenVerifier:    makeVerifier(model.RoleAdminGateway, "API Gateway"),
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard/admin?page=1&limit=20", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /dashboard/admin error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body struct {
		Status string `json:"status"`
		Data   struct {
			TrafficSummary struct {
				TotalRequests  int64   `json:"total_requests"`
				SuccessCount   int64   `json:"success_count"`
				ErrorCount     int64   `json:"error_count"`
				SuccessRatePct float64 `json:"success_rate_pct"`
				AvgDurationMS  int     `json:"avg_duration_ms"`
			} `json:"traffic_summary"`
			ServiceIndicators []struct {
				AppName string `json:"app_name"`
				Status  string `json:"status"`
			} `json:"service_indicators"`
			AuditLogs struct {
				Items []struct {
					ID        int64  `json:"id"`
					SourceApp string `json:"source_app"`
					Endpoint  string `json:"endpoint"`
					Method    string `json:"method"`
					Status    int    `json:"status"`
				} `json:"items"`
				Total int64 `json:"total"`
				Page  int   `json:"page"`
				Limit int   `json:"limit"`
			} `json:"audit_logs"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Status != "success" {
		t.Errorf("status = %q, want success", body.Status)
	}
	if body.Data.TrafficSummary.TotalRequests != 100 {
		t.Errorf("total_requests = %d, want 100", body.Data.TrafficSummary.TotalRequests)
	}
	if body.Data.TrafficSummary.SuccessRatePct != 84.0 {
		t.Errorf("success_rate_pct = %f, want 84.0", body.Data.TrafficSummary.SuccessRatePct)
	}
	if len(body.Data.ServiceIndicators) != 2 {
		t.Errorf("service_indicators count = %d, want 2", len(body.Data.ServiceIndicators))
	}
	if body.Data.ServiceIndicators[0].AppName != "Marketplace" {
		t.Errorf("service_indicators[0].app_name = %q, want Marketplace", body.Data.ServiceIndicators[0].AppName)
	}
	if body.Data.ServiceIndicators[0].Status != "inactive" {
		t.Errorf("service_indicators[0].status = %q, want inactive", body.Data.ServiceIndicators[0].Status)
	}
	if body.Data.AuditLogs.Total != 50 {
		t.Errorf("audit_logs.total = %d, want 50", body.Data.AuditLogs.Total)
	}
	if len(body.Data.AuditLogs.Items) != 1 {
		t.Errorf("audit_logs.items count = %d, want 1", len(body.Data.AuditLogs.Items))
	}
	if body.Data.AuditLogs.Items[0].SourceApp != "POS" {
		t.Errorf("audit_logs.items[0].source_app = %q, want POS", body.Data.AuditLogs.Items[0].SourceApp)
	}
	if body.Data.AuditLogs.Page != 1 {
		t.Errorf("audit_logs.page = %d, want 1", body.Data.AuditLogs.Page)
	}
	if body.Data.AuditLogs.Limit != 20 {
		t.Errorf("audit_logs.limit = %d, want 20", body.Data.AuditLogs.Limit)
	}
}

func TestDashboardAdminDefaultsPaginationParams(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		DashboardService: &stubDashboardSvc{totalLogs: 0},
		TokenVerifier:    makeVerifier(model.RoleAdminGateway, "API Gateway"),
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard/admin", nil) // tanpa page/limit
	req.Header.Set("Authorization", "Bearer any-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /dashboard/admin error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var body struct {
		Data struct {
			AuditLogs struct {
				Page  int `json:"page"`
				Limit int `json:"limit"`
			} `json:"audit_logs"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Data.AuditLogs.Page != 1 {
		t.Errorf("default page = %d, want 1", body.Data.AuditLogs.Page)
	}
	if body.Data.AuditLogs.Limit != 20 {
		t.Errorf("default limit = %d, want 20", body.Data.AuditLogs.Limit)
	}
}
