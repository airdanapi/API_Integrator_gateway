package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/notification"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

type stubNotificationSvc struct {
	listResult notification.ListResult
	listErr    error
	markResult notification.MarkReadResult
	markErr    error
	seenClaims auth.Claims
	seenPage   int
	seenLimit  int
	seenMark   notification.MarkReadRequest
}

func (s *stubNotificationSvc) List(_ context.Context, claims auth.Claims, page, limit int) (notification.ListResult, error) {
	s.seenClaims = claims
	s.seenPage = page
	s.seenLimit = limit
	return s.listResult, s.listErr
}

func (s *stubNotificationSvc) MarkRead(_ context.Context, claims auth.Claims, req notification.MarkReadRequest) (notification.MarkReadResult, error) {
	s.seenClaims = claims
	s.seenMark = req
	return s.markResult, s.markErr
}

func TestNotificationsRequireAuthentication(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{NotificationService: &stubNotificationSvc{}})
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/notifications", nil))
	if err != nil {
		t.Fatalf("GET /notifications error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", resp.StatusCode)
	}
}

func TestNotificationsListReturnsVisibleNotifications(t *testing.T) {
	svc := &stubNotificationSvc{listResult: notification.ListResult{
		Notifications: []notification.Item{{ID: 1, CreatedAt: time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC), AppName: "Marketplace", Type: model.NotificationTypeAPIInactive, Message: "inactive", IsRead: false}},
		UnreadCount:   1,
		Page:          1,
		Limit:         10,
	}}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:       makeVerifier(model.RoleAppUser, "Marketplace"),
		NotificationService: svc,
	})
	req := httptest.NewRequest(http.MethodGet, "/notifications?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /notifications error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if svc.seenClaims.AppName != "Marketplace" || svc.seenPage != 1 || svc.seenLimit != 10 {
		t.Fatalf("seen service args: claims=%#v page=%d limit=%d", svc.seenClaims, svc.seenPage, svc.seenLimit)
	}
	var body struct {
		Status string `json:"status"`
		Data   struct {
			Notifications []notification.Item `json:"notifications"`
			UnreadCount   int64               `json:"unread_count"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Status != "success" || body.Data.UnreadCount != 1 || len(body.Data.Notifications) != 1 {
		t.Fatalf("body = %#v", body)
	}
}

func TestNotificationsReadValidatesBodyAndMapsErrors(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:       makeVerifier(model.RoleAppUser, "Marketplace"),
		NotificationService: &stubNotificationSvc{},
	})
	req := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /notifications/read error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid body status = %d, want 400", resp.StatusCode)
	}

	svc := &stubNotificationSvc{markErr: notification.ErrNotFound}
	app = server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:       makeVerifier(model.RoleAppUser, "Marketplace"),
		NotificationService: svc,
	})
	req = httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBufferString(`{"notification_id":99}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("POST /notifications/read not found error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("not found status = %d, want 404", resp.StatusCode)
	}
}

func TestNotificationsReadSuccess(t *testing.T) {
	svc := &stubNotificationSvc{markResult: notification.MarkReadResult{UnreadCount: 0}}
	app := server.NewApp(config.Config{AppEnv: "test"}, server.Dependencies{
		TokenVerifier:       makeVerifier(model.RoleAdminGateway, "API Gateway"),
		NotificationService: svc,
	})
	req := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBufferString(`{"all":true}`))
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /notifications/read error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if !svc.seenMark.All {
		t.Fatalf("seenMark = %#v, want all=true", svc.seenMark)
	}
}
