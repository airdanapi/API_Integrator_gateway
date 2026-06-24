package server

import (
	"context"
	"strconv"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/dashboard"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/gofiber/fiber/v3"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

// DashboardService mendefinisikan kontrak yang dibutuhkan handler dari service layer.
type DashboardService interface {
	GetTrafficSummary(ctx context.Context, since time.Time) (dashboard.TrafficSummary, error)
	GetServiceIndicators(ctx context.Context) ([]dashboard.ServiceIndicator, error)
	GetAuditLogs(ctx context.Context, limit, offset int) ([]dashboard.AuditLogEntry, int64, error)
	GetUserDashboard(ctx context.Context, appName string, page, limit int) (dashboard.UserDashboard, error)
	GetMonitoringDashboard(ctx context.Context) (dashboard.MonitoringDashboard, error)
	GetTrafficHistory(ctx context.Context, since time.Time) ([]model.TrafficHistoryEntry, error)
}

// requireRole mengembalikan middleware yang memastikan user memiliki role tertentu.
func requireRole(allowedRoles ...model.Role) fiber.Handler {
	allowed := make(map[model.Role]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = struct{}{}
	}
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		if _, permit := allowed[claims.Role]; !permit {
			return errorResponse(c, fiber.StatusForbidden, "forbidden", "access denied for this role")
		}
		return c.Next()
	}
}

// adminDashboardHandler menangani GET /dashboard/admin.
func adminDashboardHandler(svc DashboardService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if svc == nil {
			return internalError(c)
		}

		// Pagination params
		page := parseQueryInt(c.Query("page"), defaultPage)
		limit := parseQueryInt(c.Query("limit"), defaultLimit)
		if page < 1 {
			page = defaultPage
		}
		if limit < 1 || limit > maxLimit {
			limit = defaultLimit
		}
		offset := (page - 1) * limit

		ctx := c.Context()
		since := time.Now().UTC().Add(-7 * 24 * time.Hour)

		summary, err := svc.GetTrafficSummary(ctx, since)
		if err != nil {
			return internalError(c)
		}

		indicators, err := svc.GetServiceIndicators(ctx)
		if err != nil {
			return internalError(c)
		}

		logs, total, err := svc.GetAuditLogs(ctx, limit, offset)
		if err != nil {
			return internalError(c)
		}
		
		history, err := svc.GetTrafficHistory(ctx, since)
		if err != nil {
			return internalError(c)
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"traffic_summary":    summary,
				"traffic_history":    history,
				"service_indicators": indicators,
				"audit_logs": fiber.Map{
					"items": logs,
					"total": total,
					"page":  page,
					"limit": limit,
				},
			},
		})
	}
}

// parseQueryInt mengkonversi string query param ke int; kembalikan defaultVal jika tidak valid.
func parseQueryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// userDashboardHandler menangani GET /dashboard/user.
// Memfilter data berdasarkan app_name user yang sedang login.
func userDashboardHandler(svc DashboardService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if svc == nil {
			return internalError(c)
		}

		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}

		page := parseQueryInt(c.Query("page"), defaultPage)
		limit := parseQueryInt(c.Query("limit"), defaultLimit)
		if page < 1 {
			page = defaultPage
		}
		if limit < 1 || limit > maxLimit {
			limit = defaultLimit
		}

		result, err := svc.GetUserDashboard(c.Context(), claims.AppName, page, limit)
		if err != nil {
			return internalError(c)
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   result,
		})
	}
}

// monitoringDashboardHandler menangani GET /dashboard/monitoring.
func monitoringDashboardHandler(svc DashboardService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if svc == nil {
			return internalError(c)
		}

		result, err := svc.GetMonitoringDashboard(c.Context())
		if err != nil {
			return internalError(c)
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   result,
		})
	}
}
