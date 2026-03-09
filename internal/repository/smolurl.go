package repository

import "github.com/mistic0xb/smolurl/internal/server"

type SmolURLRepository struct {
	server *server.Server
}

func NewSmolURLRepository(s *server.Server) *SmolURLRepository { return &SmolURLRepository{} }
