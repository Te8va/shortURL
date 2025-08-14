// package handler contains health check handler for database connection.
package handler

import (
	"context"
	"net/http"
)

// Pinger defines an interface for health check for database connection.
//
//go:generate mockgen -source=pinghandler.go -destination=mocks/url_pinger_mock.gen.go -package=mocks
type Pinger interface {
	PingPg(ctx context.Context) error
}

// PingHandler handles requests for health checks.
type PingHandler struct {
	pinger Pinger
}

// NewPingHandler creates new instance of PingHandler
func NewPingHandler(pinger Pinger) *PingHandler {
	return &PingHandler{pinger: pinger}
}

// PingHandler processes request to health check of database connection.
func (u *PingHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := u.pinger.PingPg(r.Context())
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
