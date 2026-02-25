package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/middleware"
	"github.com/mistic0xb/smolurl/internal/server"
)

type HealthHandler struct {
	Handler
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		Handler: NewHandler(s),
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	start := time.Now()

	logger := middleware.GetLogger(c).With().
		Str("operation", "health_check").
		Logger()

	response := map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"checks":    map[string]any{},
	}

	checks := response["checks"].(map[string]any)
	isHealthy := true

	// ---- Database check ----
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbStart := time.Now()
	if err := h.server.DB.Pool.Ping(dbCtx); err != nil {
		isHealthy = false
		checks["database"] = map[string]any{
			"status":        "unhealthy",
			"response_time": time.Since(dbStart).String(),
			"error":         err.Error(),
		}
		logger.Error().
			Err(err).
			Dur("response_time", time.Since(dbStart)).
			Msg("database health check failed")
	} else {
		checks["database"] = map[string]any{
			"status":        "healthy",
			"response_time": time.Since(dbStart).String(),
		}
		logger.Info().
			Dur("response_time", time.Since(dbStart)).
			Msg("database health check passed")
	}

	// ---- Final status ----
	if !isHealthy {
		response["status"] = "unhealthy"
		logger.Warn().
			Dur("total_duration", time.Since(start)).
			Msg("health check failed")

		return c.JSON(http.StatusServiceUnavailable, response)
	}

	logger.Info().
		Dur("total_duration", time.Since(start)).
		Msg("health check passed")

	if err := c.JSON(http.StatusOK, response); err != nil {
		logger.Error().Err(err).Msg("failed to write JSON response")
		return fmt.Errorf("failed to write JSON response: %w", err)
	}

	return nil
}
