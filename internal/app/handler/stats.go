package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
)

// URLStats defines an interface for retrieving service statistics
//
//go:generate mockgen -source=stats.go -destination=mocks/url_stats_mock.gen.go -package=mocks
type URLStats interface {
	GetStats(ctx context.Context) (urlsCount int, usersCount int, err error)
}

// StatsResponse handles response format for statistics
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// StatsHandler handles requests for service statistics
type StatsHandler struct {
	stat URLStats
	cfg  *config.Config
}

// NewStatsHandler creates a new instance of StatsHandler
func NewStatsHandler(stat URLStats, cfg *config.Config) *StatsHandler {
	return &StatsHandler{
		stat: stat,
		cfg:  cfg,
	}
}

// GetStatsHandler processes GET requests to retrieve service statistics
func (h *StatsHandler) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	urlsCount, usersCount, err := h.stat.GetStats(r.Context())
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		URLs:  urlsCount,
		Users: usersCount,
	}

	w.Header().Set(contentType, contentTypeText)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
