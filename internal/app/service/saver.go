package service

import "context"

// URLSaverServ defines the interface for a service that saves URLs
//
//go:generate mockgen -source=saver.go -destination=mocks/saver_mock.gen.go -package=mocks
type URLSaverServ interface {
	Save(ctx context.Context, userID int, url string) (string, error)
	SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error)
}

// Save delegates the save operation to repository
func (s *URLService) Save(ctx context.Context, userID int, url string) (string, error) {
	return s.saver.Save(ctx, userID, url)
}

// SaveBatch delegates the batch save operation to repository
func (s *URLService) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	return s.saver.SaveBatch(ctx, userID, urls)
}
