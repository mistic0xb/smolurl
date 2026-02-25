package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/handler"
)

func registerSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/api/status", h.Health.CheckHealth)

	r.Static("/", "static")
}
