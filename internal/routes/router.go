package routes

import (
	"net/http"

	"github.com/mistic0xb/smolurl/internal/handler"
	"github.com/mistic0xb/smolurl/internal/middleware"
	v1 "github.com/mistic0xb/smolurl/internal/routes/v1"
	"github.com/mistic0xb/smolurl/internal/server"
	"github.com/mistic0xb/smolurl/internal/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func NewRouter(s *server.Server, h *handler.Handlers, services *service.Services) *echo.Echo {
	middlewares := middleware.NewMiddleWares(s)

	router := echo.New()

	router.HTTPErrorHandler = middlewares.Global.GlobalErrorHandler

	router.Use(
		// rate limit
		echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
			Store: echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20)),
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				s.Logger.Warn().
					Str("request_id", middleware.GetRequestID(c)).
					Str("identifier", identifier).
					Str("path", c.Path()).
					Str("method", c.Request().Method).
					Str("ip", c.RealIP()).
					Msg("rate limit exceeded")

				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			},
		}),

		// check CORS
		middlewares.Global.CORS(),
		middlewares.Global.Secure(),

		// req id
		middleware.RequestID(),

		// enhance req context
		middlewares.ContextEnhancer.EnhanceContext(),

		middlewares.Global.RequestLogger(),
		middlewares.Global.Recover(),
	)

	// register system routes
	registerSystemRoutes(router, h)

	//TODO: v1 routers
	v1Routes := router.Group("/api/v1")
	v1.RegisterV1Routes(v1Routes, h, middlewares)

	return router
}
