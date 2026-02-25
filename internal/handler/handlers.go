package handler

import (
	"github.com/mistic0xb/smolurl/internal/server"
	"github.com/mistic0xb/smolurl/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	SmolURL *SmolURLHandler
}

func NewHandlers(s *server.Server, services *service.Service) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		SmolURL: NewSmolURLHandler(s, services.SmolURL),
	}
}
