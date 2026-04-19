package smolurl

type GenerateSmolURLPayload struct {
	OriginalURL    string `json:"original_url"`
	ExpirationTime int    `json:"expiration_time"`
}

type PaginatedTopSmolURLsResponse struct {
	Data []PaginatedSmolURL `json:"data"`
	Page int                `json:"page"`
}
