package service

import "context"

// URLDeleteServ defines the interface for service that deletes user URLs
//
//go:generate mockgen -source=deleter.go -destination=mocks/delete_mock.gen.go -package=mocks
type URLDeleteServ interface {
	DeleteUserURLs(ctx context.Context, ids []string, userID int) error
}

// DeleteUserURLs delegates the delete operation to repository
func (s *URLService) DeleteUserURLs(ctx context.Context, ids []string, userID int) error {
	return s.deleter.DeleteUserURLs(ctx, ids, userID)
}
