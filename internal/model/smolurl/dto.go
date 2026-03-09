package smolurl

type GenerateSmolURLPayload struct {
	OriginalURL    string `json:"original_url"`
	ExpirationTime *int   `json:"expiration_time"`
}
