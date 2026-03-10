package smolurl

import (
	"time"
)

type SmolURL struct {
	ID          uint64    `json:"id" db:"id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	SmolURL     string    `json:"smol_url" db:"smol_url"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
