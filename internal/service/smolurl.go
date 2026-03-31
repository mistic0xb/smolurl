package service

import (
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jxskiss/base62"
	"github.com/labstack/echo/v4"
	"github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/server"
)

type SmolURLService struct {
	server  *server.Server
	urlRepo *repository.SmolURLRepository
	node    *snowflake.Node
}

func NewSmolURLService(server *server.Server, urlRepo *repository.SmolURLRepository) *SmolURLService {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("failed to start snowflake node: %v\n", err)
	}
	return &SmolURLService{
		server:  server,
		urlRepo: urlRepo,
		node:    node,
	}
}

func (s *SmolURLService) GenerateSmolURL(ctx echo.Context, payload *smolurl.GenerateSmolURLPayload) (*smolurl.SmolURL, error) {
	expiresAt := time.Now().Add(30 * time.Minute)
	id := s.node.Generate()

	smolURLCode := string(base62.FormatUint(uint64(id)))

	smolURLItem, err := s.urlRepo.CreateSmolURL(ctx.Request().Context(), &smolurl.SmolURL{
		ID:          uint64(id),
		OriginalURL: payload.OriginalURL,
		SmolURL:     smolURLCode,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return smolURLItem, nil
}

func (s *SmolURLService) GetOriginalURL(ctx echo.Context, smolurlCode string) (string, error) {
	originalURL, err := s.urlRepo.GetOriginalURL(ctx.Request().Context(), smolurlCode)
	if err != nil {
		log.Fatalf("ERROR getting original url: %v", err)
	}

	return originalURL, nil
}
