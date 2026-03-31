package handler

import (
	"log"
	"net/http"

	smolurl "github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/server"
	"github.com/mistic0xb/smolurl/internal/service"

	"github.com/labstack/echo/v4"
)

type SmolURLHandler struct {
	smolURLService *service.SmolURLService
}

func NewSmolURLHandler(s *server.Server, smolURLService *service.SmolURLService) *SmolURLHandler {
	return &SmolURLHandler{
		smolURLService: smolURLService,
	}
}

func (h *SmolURLHandler) GenerateSmolURL(c echo.Context) error {
	var payload = new(smolurl.GenerateSmolURLPayload)
	// TODO: validate payload

	// bind payload
	if err := c.Bind(payload); err != nil {
		return err
	}

	// service
	res, err := h.smolURLService.GenerateSmolURL(c, payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

// takes smolURL string -> decodes it -> query the db with that id -> redirect 302
func (h *SmolURLHandler) GetUrlByID(c echo.Context) error {
	smolURLCode := c.Param("id")

	// get the original url
	originalURL, err := h.smolURLService.GetOriginalURL(c, smolURLCode)
	if err != nil {
		log.Fatalf("ERROR getting originalURL from smolURL: %v", err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, originalURL)
}
