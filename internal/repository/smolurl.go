package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/mistic0xb/smolurl/internal/model/smolurl"
	"github.com/mistic0xb/smolurl/internal/server"
)

type SmolURLRepository struct {
	server *server.Server
}

func NewSmolURLRepository(server *server.Server) *SmolURLRepository {
	return &SmolURLRepository{server: server}
}

func (r *SmolURLRepository) CreateSmolURL(ctx context.Context, payload *smolurl.SmolURL) (*smolurl.SmolURL, error) {
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
		log.Fatal(err)
		return nil, fmt.Errorf("failed to execute create smolurl query for id=%v original_url=%v: %v", payload.ID, payload.OriginalURL, err)
	}

	smolURLItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[smolurl.SmolURL])
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("failed to collect row from table:smolurls for id=%v original_url=%v: %v", payload.ID, payload.OriginalURL, err)
	}

	return &smolURLItem, nil
}

func (r *SmolURLRepository) GetOriginalURL(ctx context.Context, smolURLCode string) (string, error) {
	stmt := `SELECT original_url FROM smolurls WHERE smol_url = $1 AND expires_at > CURRENT_TIMESTAMP`
	var originalURL string
	err := r.server.DB.Pool.QueryRow(ctx, stmt, smolURLCode).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("url not found or expired")
		}
		return "", fmt.Errorf("failed to execute getOriginalURL query for smol_url=%v: %v", smolURLCode, err)
	}

	return originalURL, nil
}
