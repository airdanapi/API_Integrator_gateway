package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/dashboard"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/database"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/repository"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	startupContext, cancelStartup := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelStartup()

	db, err := database.Open(startupContext, cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(startupContext, db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	userRepository := repository.NewMySQLUserRepository(db)
	passwordHasher := auth.NewBcryptPasswordHasher()
	if cfg.SeedUsersEnabled {
		err := auth.SeedUsers(
			startupContext,
			userRepository,
			passwordHasher,
			[]auth.SeedUser{
				{
					Username: cfg.AdminSeed.Username,
					Password: cfg.AdminSeed.Password,
					Role:     model.RoleAdminGateway,
					AppName:  cfg.AdminSeed.AppName,
				},
				{
					Username: cfg.AppUserSeed.Username,
					Password: cfg.AppUserSeed.Password,
					Role:     model.RoleAppUser,
					AppName:  cfg.AppUserSeed.AppName,
				},
				{
					Username: cfg.MonitoringUserSeed.Username,
					Password: cfg.MonitoringUserSeed.Password,
					Role:     model.RoleMonitoringUser,
					AppName:  cfg.MonitoringUserSeed.AppName,
				},
			},
		)
		if err != nil {
			log.Fatalf("seed users: %v", err)
		}
	}

	jwtService := auth.NewJWTService(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTTTL,
		time.Now,
	)
	authService := auth.NewService(userRepository, passwordHasher, jwtService)

	logRepository := repository.NewMySQLLogRepository(db)
	dashboardService := dashboard.New(logRepository)

	app := server.NewApp(cfg, server.Dependencies{
		AuthService:      authService,
		TokenVerifier:    jwtService,
		DashboardService: dashboardService,
	})
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.Listen(":" + cfg.BackendPort)
	}()

	signalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("run HTTP server: %v", err)
		}
	case <-signalContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownContext); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}
}
