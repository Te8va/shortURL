package service

import (
	"context"
)

//go:generate mockgen -destination=mocks/url_saver_mock.gen.go -package=mocks . URLSaver
type URLSaver interface {
	Save(ctx context.Context, userID int, url string) (string, error)
	SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error)
}

//go:generate mockgen -destination=mocks/url_getter_mock.gen.go -package=mocks . URLGetter
type URLGetter interface {
	Get(ctx context.Context, id string) (string, bool)
	GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error)
}

//go:generate mockgen -destination=mocks/pinger_mock.gen.go -package=mocks . Pinger
type Pinger interface {
	PingPg(ctx context.Context) error
}

type URLService struct {
	saver  URLSaver
	getter URLGetter
	pinger Pinger
}

func NewURLService(saver URLSaver, getter URLGetter, pinger Pinger) *URLService {
	return &URLService{saver: saver, getter: getter, pinger: pinger}
}

func (s *URLService) PingPg(ctx context.Context) error {
	return s.pinger.PingPg(ctx)
}

func (s *URLService) Save(ctx context.Context, userID int, url string) (string, error) {
	return s.saver.Save(ctx, userID, url)
}

func (s *URLService) Get(ctx context.Context, id string) (string, bool) {
	return s.getter.Get(ctx, id)
}

func (s *URLService) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	return s.saver.SaveBatch(ctx, userID, urls)
}

func (s *URLService) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	return s.getter.GetUserURLs(ctx, userID)
}
