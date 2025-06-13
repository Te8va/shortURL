package service

import (
	"context"
)

//go:generate mockgen -source=pinger.go -destination=mocks/pinger_mock.gen.go -package=mocks
type PingerServ interface {
	PingPg(ctx context.Context) error
}

type URLService struct {
	saver   URLSaverServ
	getter  URLGetterServ
	pinger  PingerServ
	deleter URLDeleteServ
}

func NewURLService(saver URLSaverServ, getter URLGetterServ, pinger PingerServ, deleter URLDeleteServ) *URLService {
	return &URLService{saver: saver, getter: getter, pinger: pinger, deleter: deleter}
}

func (s *URLService) PingPg(ctx context.Context) error {
	return s.pinger.PingPg(ctx)
}
