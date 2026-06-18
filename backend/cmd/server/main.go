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
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	app := server.NewApp(cfg)
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
