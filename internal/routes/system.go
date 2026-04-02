package routes

import (
	"io/fs"
	"net/http"

	appstatic "github.com/mistic0xb/smolurl/static"
	"github.com/mistic0xb/smolurl/internal/handler"

	"github.com/labstack/echo/v4"
)

func registerSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/api/status", h.Health.CheckHealth)

	staticFS, _ := fs.Sub(appstatic.StaticFiles, ".")
	r.GET("/", echo.WrapHandler(http.FileServer(http.FS(staticFS))))
	r.GET("/styles.css", echo.WrapHandler(http.FileServer(http.FS(staticFS))))
}
