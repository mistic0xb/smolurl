package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/server"
)

type SmolURLService struct {
	server  *server.Server
	urlRepo *repository.SmolURLRepository
}

func NewSmolURLService(server *server.Server, urlRepo *repository.SmolURLRepository) *SmolURLService {
	return &SmolURLService{
		server:  server,
		urlRepo: urlRepo,
	}
}

func (s *SmolURLService) GenerateSmolURL(ctx echo.Context, payload *smolurl.GenerateSmolURLPayload) (*smolurl.SmolURL, error) {
	// TODO: smolURL logic

	return &smolurl.SmolURL{
		ID:             uuid.New(),
		OriginalURL:    payload.OriginalURL,
		ShortURL:       payload.OriginalURL[:7],
		ExpirationTime: time.Now().Add(1 * time.Hour),
		CreatedAt:      time.Now(),
	}, nil
}
