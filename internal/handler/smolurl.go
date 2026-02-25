package handler

import (
	"github.com/mistic0xb/smolurl/internal/server"
	"github.com/mistic0xb/smolurl/internal/service"
)

type SmolURLHandler struct {
	Handler
	smolURLService *service.SmolURLService
}

func NewSmolURLHandler(s *server.Server, smolURLService *service.SmolURLService) *SmolURLHandler {
	return &SmolURLHandler{
		Handler:        NewHandler(s),
		smolURLService: smolURLService,
	}
}
