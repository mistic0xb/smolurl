package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/handler"
)

func registerSmolURLRoutes(r *echo.Group, h *handler.SmolURLHandler) {
	// Todo operations
	urls := r.Group("/url")

	// Collection operations
	urls.POST("", h.GenerateSmolURL)

	// // Individual todo operations
	// dynamicTodo := urls.Group("/:id")
	// dynamicTodo.GET("", h.GetTodoByID)
}
