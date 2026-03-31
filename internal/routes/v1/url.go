package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/handler"
)

func registerSmolURLRoutes(r *echo.Group, h *handler.SmolURLHandler) {
	//  url operations
	urls := r.Group("/url")

	// create operations
	urls.POST("", h.GenerateSmolURL)
}
