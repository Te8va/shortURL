package domain

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type contextKey string

const UserIDKey contextKey = "userID"
