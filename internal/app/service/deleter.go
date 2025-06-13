package service

import "context"

//go:generate mockgen -source=deleter.go -destination=mocks/delete_mock.gen.go -package=mocks
type URLDeleteServ interface {
	DeleteUserURLs(ctx context.Context, ids []string, userID int) error
}

func (s *URLService) DeleteUserURLs(ctx context.Context, ids []string, userID int) error {
	return s.deleter.DeleteUserURLs(ctx, ids, userID)
}
