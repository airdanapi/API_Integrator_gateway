package notification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 50

	inactiveWindow      = 7 * 24 * time.Hour
	alertWindow         = time.Hour
	dedupWindow         = 24 * time.Hour
	minAlertSamples     = 10
	errorRateThreshold  = 0.10
	responseTimeLimitMS = 300
)

var (
	ErrInvalidRequest = errors.New("invalid notification request")
	ErrForbidden      = errors.New("notification access forbidden")
	ErrNotFound       = errors.New("notification not found")
)

var monitoredApps = []string{
	"Marketplace",
	"POS",
	"SupplierHub",
	"LogistiKita",
	"SmartBank",
}

type Store interface {
	Insert(ctx context.Context, n model.Notification) (int64, error)
	FindByID(ctx context.Context, id int64) (model.Notification, error)
	ListAll(ctx context.Context, limit, offset int) ([]model.Notification, error)
	ListByAppName(ctx context.Context, appName string, limit, offset int) ([]model.Notification, error)
	CountUnread(ctx context.Context, appName string) (int64, error)
	CountUnreadAll(ctx context.Context) (int64, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, appName string) error
	MarkAllAsReadAll(ctx context.Context) error
	ExistsRecent(ctx context.Context, appName string, notificationType model.NotificationType, since time.Time) (bool, error)
}

type LogAnalyzer interface {
	CountByStatusForApp(ctx context.Context, appName string, since time.Time) (map[int]int64, error)
	AverageDurationForApp(ctx context.Context, appName string, since time.Time) (int, int64, error)
}

type Item struct {
	ID        int64                  `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	AppName   string                 `json:"app_name"`
	Type      model.NotificationType `json:"type"`
	Message   string                 `json:"message"`
	IsRead    bool                   `json:"is_read"`
}

type ListResult struct {
	Notifications []Item `json:"notifications"`
	UnreadCount   int64  `json:"unread_count"`
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
}

type MarkReadRequest struct {
	NotificationID int64 `json:"notification_id"`
	All            bool  `json:"all"`
}

type MarkReadResult struct {
	UnreadCount int64 `json:"unread_count"`
}

type Service struct {
	store Store
	logs  LogAnalyzer
	now   func() time.Time
}

func New(store Store, logs LogAnalyzer, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{store: store, logs: logs, now: now}
}

func (s *Service) List(ctx context.Context, claims auth.Claims, page, limit int) (ListResult, error) {
	if s.store == nil {
		return ListResult{}, errors.New("notification store is nil")
	}
	page, limit = normalizePagination(page, limit)
	offset := (page - 1) * limit

	var (
		items []model.Notification
		count int64
		err   error
	)
	if seesAll(claims.Role) {
		items, err = s.store.ListAll(ctx, limit, offset)
		if err != nil {
			return ListResult{}, err
		}
		count, err = s.store.CountUnreadAll(ctx)
	} else if claims.Role == model.RoleAppUser {
		items, err = s.store.ListByAppName(ctx, claims.AppName, limit, offset)
		if err != nil {
			return ListResult{}, err
		}
		count, err = s.store.CountUnread(ctx, claims.AppName)
	} else {
		return ListResult{}, ErrForbidden
	}
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Notifications: toItems(items),
		UnreadCount:   count,
		Page:          page,
		Limit:         limit,
	}, nil
}

func (s *Service) MarkRead(ctx context.Context, claims auth.Claims, req MarkReadRequest) (MarkReadResult, error) {
	if s.store == nil {
		return MarkReadResult{}, errors.New("notification store is nil")
	}
	if req.All == (req.NotificationID > 0) {
		return MarkReadResult{}, ErrInvalidRequest
	}

	if req.All {
		if seesAll(claims.Role) {
			if err := s.store.MarkAllAsReadAll(ctx); err != nil {
				return MarkReadResult{}, err
			}
		} else if claims.Role == model.RoleAppUser {
			if err := s.store.MarkAllAsRead(ctx, claims.AppName); err != nil {
				return MarkReadResult{}, err
			}
		} else {
			return MarkReadResult{}, ErrForbidden
		}
		return s.unreadResult(ctx, claims)
	}

	item, err := s.store.FindByID(ctx, req.NotificationID)
	if errors.Is(err, sql.ErrNoRows) {
		return MarkReadResult{}, ErrNotFound
	}
	if err != nil {
		return MarkReadResult{}, err
	}
	if !visibleTo(claims, item) {
		return MarkReadResult{}, ErrNotFound
	}
	if err := s.store.MarkAsRead(ctx, req.NotificationID); err != nil {
		return MarkReadResult{}, err
	}
	return s.unreadResult(ctx, claims)
}

func (s *Service) GenerateAlerts(ctx context.Context) (int, error) {
	if s.store == nil || s.logs == nil {
		return 0, errors.New("notification dependencies are nil")
	}
	now := s.now().UTC()
	created := 0
	for _, appName := range monitoredApps {
		recentCounts, err := s.logs.CountByStatusForApp(ctx, appName, now.Add(-inactiveWindow))
		if err != nil {
			return created, err
		}
		if totalCount(recentCounts) == 0 {
			ok, err := s.insertAlertIfFresh(ctx, appName, model.NotificationTypeAPIInactive, fmt.Sprintf("API %s tidak aktif selama lebih dari 1 minggu.", appName), now)
			if err != nil {
				return created, err
			}
			if ok {
				created++
			}
		}

		windowCounts, err := s.logs.CountByStatusForApp(ctx, appName, now.Add(-alertWindow))
		if err != nil {
			return created, err
		}
		total := totalCount(windowCounts)
		if total >= minAlertSamples && errorRate(windowCounts, total) > errorRateThreshold {
			ok, err := s.insertAlertIfFresh(ctx, appName, model.NotificationTypeErrorRate, fmt.Sprintf("Error rate %s melewati 10%% dalam 1 jam terakhir.", appName), now)
			if err != nil {
				return created, err
			}
			if ok {
				created++
			}
		}

		avgMS, samples, err := s.logs.AverageDurationForApp(ctx, appName, now.Add(-alertWindow))
		if err != nil {
			return created, err
		}
		if samples >= minAlertSamples && avgMS > responseTimeLimitMS {
			ok, err := s.insertAlertIfFresh(ctx, appName, model.NotificationTypeResponseTime, fmt.Sprintf("Rata-rata response time %s melebihi 300ms dalam 1 jam terakhir.", appName), now)
			if err != nil {
				return created, err
			}
			if ok {
				created++
			}
		}
	}
	return created, nil
}

func (s *Service) insertAlertIfFresh(ctx context.Context, appName string, notificationType model.NotificationType, message string, now time.Time) (bool, error) {
	exists, err := s.store.ExistsRecent(ctx, appName, notificationType, now.Add(-dedupWindow))
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	_, err = s.store.Insert(ctx, model.Notification{
		CreatedAt: now,
		AppName:   appName,
		Type:      notificationType,
		Message:   message,
		IsRead:    false,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) unreadResult(ctx context.Context, claims auth.Claims) (MarkReadResult, error) {
	var (
		count int64
		err   error
	)
	if seesAll(claims.Role) {
		count, err = s.store.CountUnreadAll(ctx)
	} else if claims.Role == model.RoleAppUser {
		count, err = s.store.CountUnread(ctx, claims.AppName)
	} else {
		return MarkReadResult{}, ErrForbidden
	}
	if err != nil {
		return MarkReadResult{}, err
	}
	return MarkReadResult{UnreadCount: count}, nil
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 || limit > MaxLimit {
		limit = DefaultLimit
	}
	return page, limit
}

func toItems(items []model.Notification) []Item {
	result := make([]Item, 0, len(items))
	for _, item := range items {
		result = append(result, Item{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			AppName:   item.AppName,
			Type:      item.Type,
			Message:   item.Message,
			IsRead:    item.IsRead,
		})
	}
	return result
}

func seesAll(role model.Role) bool {
	return role == model.RoleAdminGateway || role == model.RoleMonitoringUser
}

func visibleTo(claims auth.Claims, item model.Notification) bool {
	return seesAll(claims.Role) || (claims.Role == model.RoleAppUser && item.AppName == claims.AppName)
}

func totalCount(counts map[int]int64) int64 {
	var total int64
	for _, count := range counts {
		total += count
	}
	return total
}

func errorRate(counts map[int]int64, total int64) float64 {
	if total == 0 {
		return 0
	}
	var failures int64
	for status, count := range counts {
		if status < 200 || status >= 300 {
			failures += count
		}
	}
	return float64(failures) / float64(total)
}
