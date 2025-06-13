package service

import "context"

//go:generate mockgen -source=getter.go -destination=mocks/getter_mock.gen.go -package=mocks
type URLGetterServ interface {
	Get(ctx context.Context, id string) (string, bool, bool)
	GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error)
}

func (s *URLService) Get(ctx context.Context, id string) (string, bool, bool) {
	return s.getter.Get(ctx, id)
}

func (s *URLService) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	return s.getter.GetUserURLs(ctx, userID)
}
