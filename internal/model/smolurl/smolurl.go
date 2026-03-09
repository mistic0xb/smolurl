package smolurl

import (
	"time"

	"github.com/google/uuid"
)

type SmolURL struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OriginalURL    string    `json:"original_url" db:"original_url"`
	ShortURL       string    `json:"short_url" db:"short_url"`
	ExpirationTime time.Time `json:"expiration_time" db:"expiration_time"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
