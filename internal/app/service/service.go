package service

import (
	"context"

	"github.com/Te8va/shortURL/internal/app/domain"
)

const length = 8

type URL struct {
	repo domain.RepositoryStore
}

func NewURL(repo domain.RepositoryStore) *URL {
	return &URL{repo: repo}
}

func (s *URL) PingPg(ctx context.Context) error {
	return s.repo.PingPg(ctx)
}

func (s *URL) Save(ctx context.Context, url string) (string, error) {
	return s.repo.Save(ctx, url)
}

func (s *URL) Get(ctx context.Context, id string) (string, bool) {
	return s.repo.Get(ctx, id)
}
