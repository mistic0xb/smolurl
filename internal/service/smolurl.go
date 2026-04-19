package service

import (
	"fmt"
	"log"
	"time"

	"github.com/mistic0xb/smolurl/internal/middleware"
	"github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/repository"
	"github.com/mistic0xb/smolurl/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

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

var (
	cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "smolurl_cache_hits_total",
		Help: "Total cache hits",
	})
	cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "smolurl_cache_misses_total",
		Help: "Total cache misses",
	})
)

func init() {
	prometheus.MustRegister(cacheHits, cacheMisses)
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
	expiresAt := time.Now().Add(time.Duration(payload.ExpirationTime) * time.Minute)
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
	reqCtx := ctx.Request().Context()
	logger := middleware.GetLogger(ctx)

	// Redis span
	reqCtx, redisSpan := otel.Tracer("smolurl").Start(reqCtx, "redis.get")
	redisSpan.SetAttributes(attribute.String("smol_url", smolurlCode))
	originalURL, err := s.server.Redis.Get(reqCtx, smolurlCode).Result()
	if err == nil {
		redisSpan.SetAttributes(attribute.Bool("cache_hit", true))
		redisSpan.End()
		cacheHits.Inc()
		logger.Debug().Str("code", smolurlCode).Msg("cache hit")
		return originalURL, nil
	}
	if err != redis.Nil {
		// real error, not just a miss
		redisSpan.RecordError(err)
		redisSpan.SetStatus(codes.Error, err.Error())
	}
	cacheMisses.Inc()
	redisSpan.SetAttributes(attribute.Bool("cache_hit", false))
	redisSpan.End()

	// DB span is inside GetOriginalURL repo call
	originalURL, err = s.urlRepo.GetOriginalURL(reqCtx, smolurlCode)
	if err != nil {
		return "", fmt.Errorf("failed to get original url: %w", err)
	}

	// Redis SET span
	reqCtx, setSpan := otel.Tracer("smolurl").Start(reqCtx, "redis.set")
	setSpan.SetAttributes(attribute.String("smol_url", smolurlCode))
	if err := s.server.Redis.Set(reqCtx, smolurlCode, originalURL, CACHE_TTL).Err(); err != nil {
		setSpan.RecordError(err)
		setSpan.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Str("code", smolurlCode).Msg("failed to cache url")
	}
	setSpan.End()

	logger.Debug().Str("code", smolurlCode).Msg("cache miss, stored in cache")
	return originalURL, nil
}

func (s *SmolURLService) GetTopURLs(ctx echo.Context, page int) (smolurl.PaginatedTopSmolURLsResponse, error) {
	var offset int = page * 10
	smolURLs, err := s.urlRepo.GetTopURL(ctx.Request().Context(), offset)
	if err != nil {
		return smolurl.PaginatedTopSmolURLsResponse{}, fmt.Errorf("failed: getTopURLs service, err:%v", err)
	}

	var paginatedResponse = smolurl.PaginatedTopSmolURLsResponse{
		Data: smolURLs,
		Page: page,
	}

	return paginatedResponse, nil
}
