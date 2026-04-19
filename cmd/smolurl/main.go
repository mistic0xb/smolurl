package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mistic0xb/smolurl/internal/config"
	"github.com/mistic0xb/smolurl/internal/handler"
	"github.com/mistic0xb/smolurl/internal/logger"
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/routes"
	"github.com/mistic0xb/smolurl/internal/server"
	"github.com/mistic0xb/smolurl/internal/service"
	"github.com/mistic0xb/smolurl/internal/telemetry"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config:" + err.Error())
	}

	logger := logger.NewLogger()

	// Init server
	srv, err := server.New(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Init repos
	repos := repository.NewRepositories(srv)

	// Init services
	services := service.NewServices(srv, repos)

	// Init handlers
	handlers := handler.NewHandlers(srv, services)

	r := routes.NewRouter(srv, handlers, services)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Tracer:
	tp, err := telemetry.NewTracerProvider(ctx, cfg.Telemetry.OTLPEndpoint)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize tracer provider")
	}
	srv.TracerProvider = tp

	// Start server
	go func() {
		err = srv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for the interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	logger.Info().Msg("server stopped properly")
}
