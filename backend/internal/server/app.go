package server

import (
	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewApp(cfg config.Config) *fiber.App {
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

	return app
}
