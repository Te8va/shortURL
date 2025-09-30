package service

import (
	"context"
)

// URLStatsServ defines the interface for a service that statistics of URLs and users
//
//go:generate mockgen -source=stats.go -destination=mocks/stats_mock.gen.go -package=mocks
type URLStatsServ interface {
	GetStats(ctx context.Context) (urlsCount int, usersCount int, err error)
}

// GetStats delegates statistiс of URLs and users
func (s *URLService) GetStats(ctx context.Context) (urlsCount int, usersCount int, err error) {
	return s.stats.GetStats(ctx)
}
