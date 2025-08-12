package service

import (
	"context"
)

// PingerServ defines the interface for a service that checks database connectivity
//
//go:generate mockgen -source=pinger.go -destination=mocks/pinger_mock.gen.go -package=mocks
type PingerServ interface {
	PingPg(ctx context.Context) error
}

// URLService is an aggregate service for working with URLs
type URLService struct {
	saver   URLSaverServ
	getter  URLGetterServ
	pinger  PingerServ
	deleter URLDeleteServ
}

// NewURLService creates a new instance of URLService with the given dependencies
func NewURLService(saver URLSaverServ, getter URLGetterServ, pinger PingerServ, deleter URLDeleteServ) *URLService {
	return &URLService{saver: saver, getter: getter, pinger: pinger, deleter: deleter}
}

// PingPg delegates the database connectivity check to repository
func (s *URLService) PingPg(ctx context.Context) error {
	return s.pinger.PingPg(ctx)
}
