package service

import (
	"context"

	"github.com/Te8va/shortURL/internal/app/domain"
)

const length = 8

type URLService struct {
	repo domain.RepositoryStore
}

func NewURLService(repo domain.RepositoryStore) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) PingPg(ctx context.Context) error {
	return s.repo.PingPg(ctx)
}

func (s *URLService) Save(ctx context.Context, url string) (string, error) {
	return s.repo.Save(ctx, url)
}

func (s *URLService) Get(ctx context.Context, id string) (string, bool) {
	return s.repo.Get(ctx, id)
}
