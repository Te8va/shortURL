package domain

import "context"


//go:generate mockgen -destination=mocks/repo_mock.gen.go -package=mocks . RepositoryStore
type RepositoryStore interface {
	Save(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, id string) (string, bool)
	PingPg(ctx context.Context) error
}

type ServiceStore interface {
	Save(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, id string) (string, bool)
	PingPg(ctx context.Context) error
}
