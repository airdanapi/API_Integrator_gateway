package server

import (
	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewApp(cfg config.Config, providedDependencies ...Dependencies) *fiber.App {
	var dependencies Dependencies
	if len(providedDependencies) > 0 {
		dependencies = providedDependencies[0]
	}

	app := fiber.New(fiber.Config{
		AppName: "API Integrator Gateway",
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		DisableColors: true,
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"service":     "api-integrator-gateway",
				"environment": cfg.AppEnv,
			},
		})
	})
	app.Get("/landing", landingHandler)
	app.Post("/auth/login", loginHandler(dependencies.AuthService))
	app.Get(
		"/auth/me",
		requireToken(dependencies.TokenVerifier),
		meHandler,
	)

	// Dashboard routes - semua memerlukan token valid.
	app.Get(
		"/dashboard/admin",
		requireToken(dependencies.TokenVerifier),
		requireRole(model.RoleAdminGateway),
		adminDashboardHandler(dependencies.DashboardService),
	)
	app.Get(
		"/dashboard/user",
		requireToken(dependencies.TokenVerifier),
		requireRole(model.RoleAppUser),
		userDashboardHandler(dependencies.DashboardService),
	)
	app.Get(
		"/dashboard/monitoring",
		requireToken(dependencies.TokenVerifier),
		requireRole(model.RoleMonitoringUser),
		monitoringDashboardHandler(dependencies.DashboardService),
	)

	app.Get(
		"/notifications",
		requireToken(dependencies.TokenVerifier),
		notificationsListHandler(dependencies.NotificationService),
	)
	app.Post(
		"/notifications/read",
		requireToken(dependencies.TokenVerifier),
		notificationReadHandler(dependencies.NotificationService),
	)

	// Chat routes - semua memerlukan token valid.
	app.Get(
		"/chat/conversations",
		requireToken(dependencies.TokenVerifier),
		chatConversationsHandler(dependencies.ChatService),
	)
	app.Get(
		"/chat/history",
		requireToken(dependencies.TokenVerifier),
		chatHistoryHandler(dependencies.ChatService),
	)
	app.Post(
		"/chat/message",
		requireToken(dependencies.TokenVerifier),
		chatMessageHandler(dependencies.ChatService),
	)
	app.Post(
		"/chat/read",
		requireToken(dependencies.TokenVerifier),
		chatReadHandler(dependencies.ChatService),
	)

	return app
}
