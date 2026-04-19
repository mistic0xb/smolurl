package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type SmolURLRepository struct {
	server *server.Server
}

func NewSmolURLRepository(server *server.Server) *SmolURLRepository {
	return &SmolURLRepository{server: server}
}

func (r *SmolURLRepository) CreateSmolURL(ctx context.Context, payload *smolurl.SmolURL) (*smolurl.SmolURL, error) {
	// start a child span
	ctx, span := otel.Tracer("smolurl").Start(ctx, "db.CreateSmolURL")
	defer span.End()
	span.SetAttributes(attribute.String("original_url", payload.OriginalURL))

	stmt := `
		INSERT INTO
			smolurls (
				id,
				original_url,
				smol_url,
				expires_at,
				created_at
			)
		VALUES
			(
				@id,
				@original_url,
				@smol_url,
				@expires_at,
				@created_at
			)
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":           payload.ID,
		"original_url": payload.OriginalURL,
		"smol_url":     payload.SmolURL,
		"expires_at":   payload.ExpiresAt,
		"created_at":   payload.CreatedAt,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to execute create smolurl query for id=%v original_url=%v: %v", payload.ID, payload.OriginalURL, err)
	}

	smolURLItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[smolurl.SmolURL])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to collect row from table:smolurls for id=%v original_url=%v: %v", payload.ID, payload.OriginalURL, err)
	}

	return &smolURLItem, nil
}

func (r *SmolURLRepository) GetOriginalURL(ctx context.Context, smolURLCode string) (string, error) {
	// start a child span
	ctx, span := otel.Tracer("smolurl").Start(ctx, "db.GetOriginalURL")
	defer span.End()
	span.SetAttributes(attribute.String("smol_url", smolURLCode))

	stmt := `SELECT original_url FROM smolurls WHERE smol_url = $1 AND expires_at > NOW()`
	var originalURL string
	err := r.server.DB.Pool.QueryRow(ctx, stmt, smolURLCode).Scan(&originalURL)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("url not found or expired")
		}
		return "", fmt.Errorf("failed to execute getOriginalURL query for smol_url=%v: %v", smolURLCode, err)
	}

	return originalURL, nil
}

func (r *SmolURLRepository) GetTopURL(ctx context.Context, offset int) ([]smolurl.PaginatedSmolURL, error) {
	ctx, span := otel.Tracer("smolurl").Start(ctx, "db.GetTopURL")
	defer span.End()
	span.SetAttributes(attribute.Int("offset", offset))

	stmt := `
		SELECT original_url,smol_url,created_at FROM smolurls
		WHERE expires_at > now()
		ORDER BY created_at DESC
		LIMIT 10 OFFSET $1
		`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to execute getTopURL")
	}

	smolURLs, err := pgx.CollectRows(rows, pgx.RowToStructByName[smolurl.PaginatedSmolURL])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return []smolurl.PaginatedSmolURL{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:smolurls for offset=%v: %w", offset, err)
	}

	return smolURLs, nil
}
