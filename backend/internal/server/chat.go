package server

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/chat"
	"github.com/gofiber/fiber/v3"
)

type ChatService interface {
	ListConversations(ctx context.Context, claims auth.Claims) (chat.ConversationListResult, error)
	History(ctx context.Context, claims auth.Claims, conversationID string, page, limit int) (chat.HistoryResult, error)
	SendMessage(ctx context.Context, claims auth.Claims, request chat.SendMessageRequest) (chat.SendMessageResult, error)
	MarkRead(ctx context.Context, claims auth.Claims, request chat.MarkReadRequest) (chat.MarkReadResult, error)
}

func chatConversationsHandler(service ChatService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		result, err := service.ListConversations(c.Context(), claims)
		if err != nil {
			return chatError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": result})
	}
}

func chatHistoryHandler(service ChatService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		conversationID := strings.TrimSpace(c.Query("conversation_id"))
		if conversationID == "" {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "conversation_id is required")
		}
		page := parseQueryInt(c.Query("page"), chat.DefaultPage)
		limit := parseQueryInt(c.Query("limit"), chat.DefaultLimit)
		result, err := service.History(c.Context(), claims, conversationID, page, limit)
		if err != nil {
			return chatError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": result})
	}
}

func chatMessageHandler(service ChatService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var request chat.SendMessageRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		if strings.TrimSpace(request.Message) == "" {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "message is required")
		}
		result, err := service.SendMessage(c.Context(), claims, request)
		if err != nil {
			return chatError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": result})
	}
}

func chatReadHandler(service ChatService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var request chat.MarkReadRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		if strings.TrimSpace(request.ConversationID) == "" {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "conversation_id is required")
		}
		result, err := service.MarkRead(c.Context(), claims, request)
		if err != nil {
			return chatError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": result})
	}
}

func chatError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, chat.ErrInvalidRequest):
		return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "invalid chat request")
	case errors.Is(err, chat.ErrForbidden):
		return errorResponse(c, fiber.StatusForbidden, "forbidden", "access denied for this role")
	case errors.Is(err, chat.ErrNotFound):
		return errorResponse(c, fiber.StatusNotFound, "not_found", "chat conversation not found")
	default:
		return internalError(c)
	}
}
