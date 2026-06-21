package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/gateway"
	"github.com/gofiber/fiber/v3"
)

// GatewayService mendefinisikan kontrak yang dibutuhkan handler dari gateway service layer.
type GatewayService interface {
	ForwardPayment(ctx context.Context, claims auth.Claims, req gateway.PaymentRequest) (gateway.GatewayResponse, error)
	ForwardSmartBank(ctx context.Context, claims auth.Claims, req gateway.SmartBankRequest) (gateway.GatewayResponse, error)
	ForwardMarketplace(ctx context.Context, claims auth.Claims, req gateway.MarketplaceRequest) (gateway.GatewayResponse, error)
	ForwardLogistics(ctx context.Context, claims auth.Claims, req gateway.LogisticsRequest) (gateway.GatewayResponse, error)
	ForwardSupplier(ctx context.Context, claims auth.Claims, req gateway.SupplierRequest) (gateway.GatewayResponse, error)
}

func gatewayPaymentHandler(service GatewayService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var req gateway.PaymentRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		resp, err := service.ForwardPayment(c.Context(), claims, req)
		if err != nil {
			return gatewayError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func gatewaySmartBankHandler(service GatewayService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var req gateway.SmartBankRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		resp, err := service.ForwardSmartBank(c.Context(), claims, req)
		if err != nil {
			return gatewayError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func gatewayMarketplaceHandler(service GatewayService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var req gateway.MarketplaceRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		resp, err := service.ForwardMarketplace(c.Context(), claims, req)
		if err != nil {
			return gatewayError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func gatewayLogisticsHandler(service GatewayService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var req gateway.LogisticsRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		resp, err := service.ForwardLogistics(c.Context(), claims, req)
		if err != nil {
			return gatewayError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func gatewaySupplierHandler(service GatewayService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}
		claims, ok := c.Locals("auth_claims").(auth.Claims)
		if !ok {
			return unauthorized(c)
		}
		var req gateway.SupplierRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		}
		resp, err := service.ForwardSupplier(c.Context(), claims, req)
		if err != nil {
			return gatewayError(c, err)
		}
		return c.JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func gatewayError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidPayload):
		return errorResponse(c, fiber.StatusBadRequest, "invalid_payload", err.Error())
	case errors.Is(err, gateway.ErrForbidden):
		return errorResponse(c, fiber.StatusForbidden, "forbidden", "access denied for this role")
	case errors.Is(err, gateway.ErrUpstreamNotConfigured):
		return errorResponse(c, fiber.StatusServiceUnavailable, "upstream_not_configured", "upstream service not configured")
	case errors.Is(err, gateway.ErrUpstreamFailure):
		return errorResponse(c, fiber.StatusBadGateway, "upstream_failure", "upstream service failure")
	default:
		return internalError(c)
	}
}
