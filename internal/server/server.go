package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mistic0xb/smolurl/internal/config"
	"github.com/mistic0xb/smolurl/internal/database"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Server struct {
	Config     *config.Config
	Logger     *zerolog.Logger
	DB         *database.Database
	Redis      *redis.Client
	httpServer *http.Server
}

func New(cfg *config.Config, logger *zerolog.Logger) (*Server, error) {
	// Init Database
	db, err := database.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Test Redis connection
	redisCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(redisCtx).Err(); err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Redis, continuing without Redis")
	}

	// Run migrations before anything else
	migrationCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.RunMigrations(migrationCtx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	server := &Server{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Redis:  redisClient,
	}

	return server, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}
	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}
