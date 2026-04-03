package service

import (
	"fmt"
	"log"
	"time"

	"github.com/mistic0xb/smolurl/internal/middleware"
	"github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/server"

	"github.com/bwmarrin/snowflake"
	"github.com/jxskiss/base62"
	"github.com/labstack/echo/v4"
)

const CACHE_TTL = 5 * time.Minute

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
	logger := middleware.GetLogger(ctx)

	originalURL, err := s.server.Redis.Get(ctx.Request().Context(), smolurlCode).Result()
	if err == nil {
		logger.Debug().Str("code", smolurlCode).Msg("cache hit")
		return originalURL, nil
	}

	originalURL, err = s.urlRepo.GetOriginalURL(ctx.Request().Context(), smolurlCode)
	if err != nil {
		return "", fmt.Errorf("failed to get original url: %w", err)
	}

	if err := s.server.Redis.Set(ctx.Request().Context(), smolurlCode, originalURL, CACHE_TTL).Err(); err != nil {
		logger.Warn().Err(err).Str("code", smolurlCode).Msg("failed to cache url")
	}
	logger.Debug().Str("code", smolurlCode).Msg("cache miss, stored in cache")

	return originalURL, nil
}
