package server

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/gofiber/fiber/v3"
)

type LoginService interface {
	Login(context.Context, auth.LoginRequest) (auth.LoginResult, error)
}

type TokenVerifier interface {
	Validate(token string) (auth.Claims, error)
}

type Dependencies struct {
	AuthService      LoginService
	TokenVerifier    TokenVerifier
	DashboardService DashboardService
}

func loginHandler(service LoginService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if service == nil {
			return internalError(c)
		}

		var request auth.LoginRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return errorResponse(
				c,
				fiber.StatusBadRequest,
				"invalid_request",
				"request body must be valid JSON",
			)
		}
		request.Username = strings.TrimSpace(request.Username)
		request.AppName = strings.TrimSpace(request.AppName)
		if request.Username == "" ||
			request.Password == "" ||
			request.AppName == "" {
			return errorResponse(
				c,
				fiber.StatusBadRequest,
				"invalid_request",
				"username, password, and app_name are required",
			)
		}

		result, err := service.Login(c.Context(), request)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return errorResponse(
				c,
				fiber.StatusUnauthorized,
				"invalid_credentials",
				"invalid credentials",
			)
		}
		if err != nil {
			return internalError(c)
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   result,
		})
	}
}

func requireToken(verifier TokenVerifier) fiber.Handler {
	return func(c fiber.Ctx) error {
		authorization := strings.TrimSpace(c.Get(fiber.HeaderAuthorization))
		parts := strings.Fields(authorization)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			parts[1] == "" ||
			verifier == nil {
			return unauthorized(c)
		}

		claims, err := verifier.Validate(parts[1])
		if err != nil {
			return unauthorized(c)
		}
		c.Locals("auth_claims", claims)
		return c.Next()
	}
}

func meHandler(c fiber.Ctx) error {
	claims, ok := c.Locals("auth_claims").(auth.Claims)
	if !ok {
		return unauthorized(c)
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id":  claims.Subject,
			"username": claims.Username,
			"role":     claims.Role,
			"app_name": claims.AppName,
		},
	})
}

func unauthorized(c fiber.Ctx) error {
	return errorResponse(
		c,
		fiber.StatusUnauthorized,
		"unauthorized",
		"authentication is required",
	)
}

func internalError(c fiber.Ctx) error {
	return errorResponse(
		c,
		fiber.StatusInternalServerError,
		"internal_error",
		"internal server error",
	)
}

func errorResponse(
	c fiber.Ctx,
	status int,
	code string,
	message string,
) error {
	return c.Status(status).JSON(fiber.Map{
		"status": "error",
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
