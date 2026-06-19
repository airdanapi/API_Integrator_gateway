package dashboard_test

import (
	"context"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/dashboard"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

// ─── stub LogQuerier ────────────────────────────────────────────────────────

type stubLogQuerier struct {
	byStatus      map[int]int64
	byStatusForApp map[int]int64 // returned by CountByStatusForApp
	byApp         map[string]int64
	logs          []model.RequestLog
}

func (s *stubLogQuerier) ListRecent(_ context.Context, limit, offset int) ([]model.RequestLog, error) {
	end := offset + limit
	if end > len(s.logs) {
		end = len(s.logs)
	}
	if offset >= len(s.logs) {
		return nil, nil
	}
	return s.logs[offset:end], nil
}

func (s *stubLogQuerier) ListBySourceApp(_ context.Context, _ string, limit, offset int) ([]model.RequestLog, error) {
	end := offset + limit
	if end > len(s.logs) {
		end = len(s.logs)
	}
	if offset >= len(s.logs) {
		return nil, nil
	}
	return s.logs[offset:end], nil
}

func (s *stubLogQuerier) CountByStatus(_ context.Context, _ time.Time) (map[int]int64, error) {
	return s.byStatus, nil
}

func (s *stubLogQuerier) CountByStatusForApp(_ context.Context, _ string, _ time.Time) (map[int]int64, error) {
	if s.byStatusForApp != nil {
		return s.byStatusForApp, nil
	}
	return s.byStatus, nil
}

func (s *stubLogQuerier) CountBySourceApp(_ context.Context, _ time.Time) (map[string]int64, error) {
	return s.byApp, nil
}

// ─── GetTrafficSummary ──────────────────────────────────────────────────────

func TestGetTrafficSummary_CalculatesRateCorrectly(t *testing.T) {
	querier := &stubLogQuerier{
		byStatus: map[int]int64{
			200: 80,
			400: 10,
			500: 10,
		},
	}
	svc := dashboard.New(querier)
	summary, err := svc.GetTrafficSummary(context.Background(), time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("GetTrafficSummary() error: %v", err)
	}
	if summary.TotalRequests != 100 {
		t.Errorf("TotalRequests = %d, want 100", summary.TotalRequests)
	}
	if summary.SuccessCount != 80 {
		t.Errorf("SuccessCount = %d, want 80", summary.SuccessCount)
	}
	if summary.ErrorCount != 20 {
		t.Errorf("ErrorCount = %d, want 20", summary.ErrorCount)
	}
	if summary.SuccessRatePct != 80.0 {
		t.Errorf("SuccessRatePct = %f, want 80.0", summary.SuccessRatePct)
	}
}

func TestGetTrafficSummary_EmptyReturnsZeroRate(t *testing.T) {
	querier := &stubLogQuerier{byStatus: map[int]int64{}}
	svc := dashboard.New(querier)
	summary, err := svc.GetTrafficSummary(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("GetTrafficSummary() error: %v", err)
	}
	if summary.TotalRequests != 0 {
		t.Errorf("TotalRequests = %d, want 0", summary.TotalRequests)
	}
	if summary.SuccessRatePct != 0.0 {
		t.Errorf("SuccessRatePct = %f, want 0.0", summary.SuccessRatePct)
	}
}

// ─── GetServiceIndicators ───────────────────────────────────────────────────

func TestGetServiceIndicators_ActiveAndInactive(t *testing.T) {
	querier := &stubLogQuerier{
		byApp: map[string]int64{
			"Marketplace": 10,
			"POS":         5,
			// SupplierHub, LogistiKita, SmartBank: tidak ada → inactive
		},
	}
	svc := dashboard.New(querier)
	indicators, err := svc.GetServiceIndicators(context.Background())
	if err != nil {
		t.Fatalf("GetServiceIndicators() error: %v", err)
	}
	if len(indicators) != 5 {
		t.Fatalf("indicators count = %d, want 5", len(indicators))
	}

	byName := make(map[string]dashboard.ServiceIndicator)
	for _, ind := range indicators {
		byName[ind.AppName] = ind
	}

	if byName["Marketplace"].Status != "active" {
		t.Errorf("Marketplace status = %q, want active", byName["Marketplace"].Status)
	}
	if byName["POS"].Status != "active" {
		t.Errorf("POS status = %q, want active", byName["POS"].Status)
	}
	if byName["SupplierHub"].Status != "inactive" {
		t.Errorf("SupplierHub status = %q, want inactive", byName["SupplierHub"].Status)
	}
	if byName["LogistiKita"].Status != "inactive" {
		t.Errorf("LogistiKita status = %q, want inactive", byName["LogistiKita"].Status)
	}
	if byName["SmartBank"].Status != "inactive" {
		t.Errorf("SmartBank status = %q, want inactive", byName["SmartBank"].Status)
	}
}

func TestGetServiceIndicators_AllInactive(t *testing.T) {
	querier := &stubLogQuerier{byApp: map[string]int64{}}
	svc := dashboard.New(querier)
	indicators, err := svc.GetServiceIndicators(context.Background())
	if err != nil {
		t.Fatalf("GetServiceIndicators() error: %v", err)
	}
	for _, ind := range indicators {
		if ind.Status != "inactive" {
			t.Errorf("%s status = %q, want inactive", ind.AppName, ind.Status)
		}
	}
}

// ─── GetAuditLogs ───────────────────────────────────────────────────────────

func TestGetAuditLogs_PaginatesCorrectly(t *testing.T) {
	now := time.Now().UTC()
	dur := 100
	logs := []model.RequestLog{
		{ID: 1, SourceApp: "Marketplace", Endpoint: "/gateway/payment", Method: "POST", Status: 200, Timestamp: now, DurationMS: &dur},
		{ID: 2, SourceApp: "POS", Endpoint: "/gateway/payment", Method: "POST", Status: 400, Timestamp: now},
		{ID: 3, SourceApp: "SmartBank", Endpoint: "/gateway/smartbank", Method: "POST", Status: 200, Timestamp: now},
	}
	querier := &stubLogQuerier{
		logs:     logs,
		byStatus: map[int]int64{200: 2, 400: 1},
	}
	svc := dashboard.New(querier)

	// Page 1, limit 2
	entries, total, err := svc.GetAuditLogs(context.Background(), 2, 0)
	if err != nil {
		t.Fatalf("GetAuditLogs() error: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(entries) != 2 {
		t.Errorf("entries count = %d, want 2", len(entries))
	}
	if entries[0].ID != 1 {
		t.Errorf("entries[0].ID = %d, want 1", entries[0].ID)
	}
	if entries[0].DurationMS == nil || *entries[0].DurationMS != 100 {
		t.Errorf("entries[0].DurationMS not set correctly")
	}
	if entries[1].DurationMS != nil {
		t.Errorf("entries[1].DurationMS should be nil")
	}
}

// ─── GetUserDashboard ────────────────────────────────────────────────────────

func TestGetUserDashboard_ActiveApp(t *testing.T) {
	now := time.Now().UTC()
	dur := 120
	logs := []model.RequestLog{
		{ID: 10, SourceApp: "Marketplace", Endpoint: "/gateway/payment", Method: "POST", Status: 200, Timestamp: now, DurationMS: &dur},
		{ID: 11, SourceApp: "Marketplace", Endpoint: "/gateway/payment", Method: "POST", Status: 400, Timestamp: now},
	}
	querier := &stubLogQuerier{
		byStatusForApp: map[int]int64{200: 10, 400: 2}, // recent counts
		byStatus:       map[int]int64{200: 10, 400: 2}, // fallback / all-time
		logs:           logs,
	}
	svc := dashboard.New(querier)
	result, err := svc.GetUserDashboard(context.Background(), "Marketplace", 1, 20)
	if err != nil {
		t.Fatalf("GetUserDashboard() error: %v", err)
	}
	if result.MyApp != "Marketplace" {
		t.Errorf("MyApp = %q, want Marketplace", result.MyApp)
	}
	if result.ServiceStatus != "active" {
		t.Errorf("ServiceStatus = %q, want active", result.ServiceStatus)
	}
	if result.TrafficSummary.TotalRequests != 12 {
		t.Errorf("TotalRequests = %d, want 12", result.TrafficSummary.TotalRequests)
	}
	if result.TrafficSummary.SuccessCount != 10 {
		t.Errorf("SuccessCount = %d, want 10", result.TrafficSummary.SuccessCount)
	}
	if result.Page != 1 {
		t.Errorf("Page = %d, want 1", result.Page)
	}
	if result.Limit != 20 {
		t.Errorf("Limit = %d, want 20", result.Limit)
	}
	if len(result.RecentLogs) != 2 {
		t.Errorf("RecentLogs count = %d, want 2", len(result.RecentLogs))
	}
}

func TestGetUserDashboard_InactiveApp(t *testing.T) {
	querier := &stubLogQuerier{
		byStatusForApp: map[int]int64{}, // no recent requests
		byStatus:       map[int]int64{},
		logs:           nil,
	}
	svc := dashboard.New(querier)
	result, err := svc.GetUserDashboard(context.Background(), "LogistiKita", 1, 20)
	if err != nil {
		t.Fatalf("GetUserDashboard() error: %v", err)
	}
	if result.ServiceStatus != "inactive" {
		t.Errorf("ServiceStatus = %q, want inactive", result.ServiceStatus)
	}
	if result.TrafficSummary.TotalRequests != 0 {
		t.Errorf("TotalRequests = %d, want 0", result.TrafficSummary.TotalRequests)
	}
	if len(result.RecentLogs) != 0 {
		t.Errorf("RecentLogs should be empty")
	}
}

// ─── GetMonitoringDashboard ──────────────────────────────────────────────────

func TestGetMonitoringDashboard_ReturnsAllComponents(t *testing.T) {
	querier := &stubLogQuerier{
		byStatus: map[int]int64{200: 50, 400: 10, 500: 5},
		byApp: map[string]int64{
			"Marketplace": 20,
			"POS":         15,
			"SupplierHub": 10,
			"LogistiKita": 8,
			"SmartBank":   12,
		},
	}
	svc := dashboard.New(querier)
	result, err := svc.GetMonitoringDashboard(context.Background())
	if err != nil {
		t.Fatalf("GetMonitoringDashboard() error: %v", err)
	}

	// Traffic summary overall
	if result.TrafficSummary.TotalRequests != 65 {
		t.Errorf("TotalRequests = %d, want 65", result.TrafficSummary.TotalRequests)
	}
	if result.TrafficSummary.SuccessCount != 50 {
		t.Errorf("SuccessCount = %d, want 50", result.TrafficSummary.SuccessCount)
	}

	// Service indicators: semua 5 app harus muncul
	if len(result.ServiceIndicators) != 5 {
		t.Errorf("ServiceIndicators count = %d, want 5", len(result.ServiceIndicators))
	}

	// App breakdown: semua 5 known apps harus ada
	if len(result.AppBreakdown) != 5 {
		t.Fatalf("AppBreakdown count = %d, want 5", len(result.AppBreakdown))
	}
	byName := make(map[string]dashboard.AppStat)
	for _, stat := range result.AppBreakdown {
		byName[stat.AppName] = stat
	}
	if byName["Marketplace"].TotalRequests != 20 {
		t.Errorf("Marketplace total = %d, want 20", byName["Marketplace"].TotalRequests)
	}
	if byName["SmartBank"].TotalRequests != 12 {
		t.Errorf("SmartBank total = %d, want 12", byName["SmartBank"].TotalRequests)
	}
}
