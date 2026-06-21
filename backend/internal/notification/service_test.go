package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

type stubStore struct {
	allItems     []model.Notification
	appItems     []model.Notification
	byID         model.Notification
	findErr      error
	unreadAll    int64
	unreadByApp  int64
	inserted     []model.Notification
	recentByKey  map[string]bool
	markedID     int64
	markedAll    bool
	markedAllApp string
}

func (s *stubStore) Insert(_ context.Context, n model.Notification) (int64, error) {
	s.inserted = append(s.inserted, n)
	return int64(len(s.inserted)), nil
}
func (s *stubStore) ListAll(_ context.Context, _, _ int) ([]model.Notification, error) {
	return s.allItems, nil
}
func (s *stubStore) ListByAppName(_ context.Context, _ string, _, _ int) ([]model.Notification, error) {
	return s.appItems, nil
}
func (s *stubStore) CountUnreadAll(_ context.Context) (int64, error) { return s.unreadAll, nil }
func (s *stubStore) CountUnread(_ context.Context, _ string) (int64, error) {
	return s.unreadByApp, nil
}
func (s *stubStore) FindByID(_ context.Context, _ int64) (model.Notification, error) {
	return s.byID, s.findErr
}
func (s *stubStore) MarkAsRead(_ context.Context, id int64) error { s.markedID = id; return nil }
func (s *stubStore) MarkAllAsReadAll(_ context.Context) error     { s.markedAll = true; return nil }
func (s *stubStore) MarkAllAsRead(_ context.Context, appName string) error {
	s.markedAllApp = appName
	return nil
}
func (s *stubStore) ExistsRecent(_ context.Context, appName string, typ model.NotificationType, _ time.Time) (bool, error) {
	return s.recentByKey[appName+"|"+string(typ)], nil
}

type stubLogs struct {
	statusByApp map[string]map[int]int64
	avgByApp    map[string]int
	avgCount    map[string]int64
}

func (s stubLogs) CountByStatusForApp(_ context.Context, appName string, _ time.Time) (map[int]int64, error) {
	if s.statusByApp[appName] == nil {
		return map[int]int64{}, nil
	}
	return s.statusByApp[appName], nil
}
func (s stubLogs) AverageDurationForApp(_ context.Context, appName string, _ time.Time) (int, int64, error) {
	return s.avgByApp[appName], s.avgCount[appName], nil
}

func claims(role model.Role, appName string) auth.Claims {
	return auth.Claims{Username: "tester", Role: role, AppName: appName}
}

func TestServiceListUsesRoleVisibility(t *testing.T) {
	store := &stubStore{
		allItems:    []model.Notification{{ID: 1, AppName: "Marketplace"}, {ID: 2, AppName: "POS"}},
		appItems:    []model.Notification{{ID: 3, AppName: "Marketplace"}},
		unreadAll:   2,
		unreadByApp: 1,
	}
	svc := New(store, stubLogs{}, func() time.Time { return time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC) })

	adminResult, err := svc.List(context.Background(), claims(model.RoleAdminGateway, "API Gateway"), 1, 10)
	if err != nil {
		t.Fatalf("admin List() error: %v", err)
	}
	if len(adminResult.Notifications) != 2 || adminResult.UnreadCount != 2 {
		t.Fatalf("admin List() = %#v", adminResult)
	}

	appResult, err := svc.List(context.Background(), claims(model.RoleAppUser, "Marketplace"), 1, 10)
	if err != nil {
		t.Fatalf("app List() error: %v", err)
	}
	if len(appResult.Notifications) != 1 || appResult.Notifications[0].AppName != "Marketplace" || appResult.UnreadCount != 1 {
		t.Fatalf("app List() = %#v", appResult)
	}
}

func TestServiceMarkReadRejectsInvisibleNotification(t *testing.T) {
	store := &stubStore{byID: model.Notification{ID: 9, AppName: "POS"}}
	svc := New(store, stubLogs{}, time.Now)

	_, err := svc.MarkRead(context.Background(), claims(model.RoleAppUser, "Marketplace"), MarkReadRequest{NotificationID: 9})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("MarkRead() error = %v, want ErrNotFound", err)
	}
	if store.markedID != 0 {
		t.Fatalf("markedID = %d, want 0", store.markedID)
	}
}

func TestServiceMarkAllUsesVisibleScope(t *testing.T) {
	store := &stubStore{unreadByApp: 0}
	svc := New(store, stubLogs{}, time.Now)

	result, err := svc.MarkRead(context.Background(), claims(model.RoleAppUser, "Marketplace"), MarkReadRequest{All: true})
	if err != nil {
		t.Fatalf("MarkRead(all) error: %v", err)
	}
	if store.markedAllApp != "Marketplace" || store.markedAll {
		t.Fatalf("mark all scope = all:%v app:%q", store.markedAll, store.markedAllApp)
	}
	if result.UnreadCount != 0 {
		t.Fatalf("UnreadCount = %d, want 0", result.UnreadCount)
	}
}

func TestGenerateAlertsCreatesConfiguredAlertTypesAndDedups(t *testing.T) {
	now := time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC)
	store := &stubStore{recentByKey: map[string]bool{
		"POS|error_rate": true,
	}}
	logs := stubLogs{
		statusByApp: map[string]map[int]int64{
			"Marketplace": {},
			"POS":         {200: 8, 500: 2},
			"SupplierHub": {200: 10},
		},
		avgByApp: map[string]int{"SupplierHub": 350},
		avgCount: map[string]int64{"SupplierHub": 10},
	}
	svc := New(store, logs, func() time.Time { return now })

	created, err := svc.GenerateAlerts(context.Background())
	if err != nil {
		t.Fatalf("GenerateAlerts() error: %v", err)
	}
	if created == 0 {
		t.Fatal("GenerateAlerts() created 0 alerts, want at least inactive and response_time")
	}
	if hasInserted(store.inserted, "POS", model.NotificationTypeErrorRate) {
		t.Fatal("GenerateAlerts() inserted deduped POS error_rate alert")
	}
	if !hasInserted(store.inserted, "Marketplace", model.NotificationTypeAPIInactive) {
		t.Fatal("GenerateAlerts() did not insert Marketplace api_inactive alert")
	}
	if !hasInserted(store.inserted, "SupplierHub", model.NotificationTypeResponseTime) {
		t.Fatal("GenerateAlerts() did not insert SupplierHub response_time alert")
	}
}

func hasInserted(items []model.Notification, appName string, typ model.NotificationType) bool {
	for _, item := range items {
		if item.AppName == appName && item.Type == typ {
			return true
		}
	}
	return false
}
