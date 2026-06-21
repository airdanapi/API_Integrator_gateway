package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/notification"
	"github.com/gofiber/fiber/v3"
)

const (
	defaultNotificationPage  = notification.DefaultPage
	defaultNotificationLimit = notification.DefaultLimit
	maxNotificationLimit     = notification.MaxLimit
)

type NotificationService interface {
	List(ctx context.Context, claims auth.Claims, page, limit int) (notification.ListResult, error)
	MarkRead(ctx context.Context, claims auth.Claims, request notification.MarkReadRequest) (notification.MarkReadResult, error)
}

func notificationsListHandler(service NotificationService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}

		page := parseQueryInt(c.Query("page"), defaultNotificationPage)
		limit := parseQueryInt(c.Query("limit"), defaultNotificationLimit)
		if page < 1 {
			page = defaultNotificationPage
		}
		if limit < 1 || limit > maxNotificationLimit {
			limit = defaultNotificationLimit
		}

		result, err := service.List(c.Context(), claims, page, limit)
		if err != nil {
			return notificationError(c, err)
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   result,
		})
	}
}

func notificationReadHandler(service NotificationService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}

		var request notification.MarkReadRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		if request.All == (request.NotificationID > 0) {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "provide either notification_id or all=true")
		}

		result, err := service.MarkRead(c.Context(), claims, request)
		if err != nil {
			return notificationError(c, err)
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   result,
		})
	}
}

func notificationError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, notification.ErrInvalidRequest):
		return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "invalid notification request")
	case errors.Is(err, notification.ErrForbidden):
		return errorResponse(c, fiber.StatusForbidden, "forbidden", "access denied for this role")
	case errors.Is(err, notification.ErrNotFound):
		return errorResponse(c, fiber.StatusNotFound, "not_found", "notification not found")
	default:
		return internalError(c)
	}
}
