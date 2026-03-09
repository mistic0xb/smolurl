package service

import (
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/server"
)

type Services struct {
	SmolURL *SmolURLService
}

func NewServices(s *server.Server, repos *repository.Repositories) *Services {
	return &Services{
		SmolURL: NewSmolURLService(s, repos.SmolURL),
	}
}
