package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/handler"
	"github.com/mistic0xb/smolurl/internal/middleware"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {
	// Register url routes
	registerSmolURLRoutes(router, handlers.SmolURL)
}
