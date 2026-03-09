package repository

import "github.com/mistic0xb/smolurl/internal/server"

type Repositories struct {
	SmolURL *SmolURLRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		SmolURL: NewSmolURLRepository(s),
	}
}
