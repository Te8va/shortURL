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

// UserIDKey is the key used to store the user ID in context
const UserIDKey contextKey = "userID"
