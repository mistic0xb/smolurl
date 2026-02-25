package service

import (
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/server"
)

type SmolURLService struct {
	server  *server.Server
	urlRepo *repository.SmolURLRepository
}
