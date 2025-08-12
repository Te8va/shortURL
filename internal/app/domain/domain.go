package domain

// ShortenRequest represents request to URL.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents response containing userID .
type ShortenResponse struct {
	Result string `json:"result"`
}

type contextKey string

const UserIDKey contextKey = "userID"
