package service

import "context"

// URLGetterServ defines the interface for a service that retrieves URLs
//
//go:generate mockgen -source=getter.go -destination=mocks/getter_mock.gen.go -package=mocks
type URLGetterServ interface {
	Get(ctx context.Context, id string) (string, bool, bool)
	GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error)
}

// Get delegates the retrieval operation to repository
func (s *URLService) Get(ctx context.Context, id string) (string, bool, bool) {
	return s.getter.Get(ctx, id)
}

// GetUserURLs delegates the retrieval of the user's URLs to repository
func (s *URLService) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	return s.getter.GetUserURLs(ctx, userID)
}
