package handler

import (
	"context"
	"net/http"
)

//go:generate mockgen -source=pinghandler.go -destination=mocks/url_pinger_mock.gen.go -package=mocks
type Pinger interface {
	PingPg(ctx context.Context) error
}

type PingHandler struct {
	pinger Pinger
}

func NewPingHandler(pinger Pinger) *PingHandler {
	return &PingHandler{pinger: pinger}
}

func (u *PingHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := u.pinger.PingPg(r.Context())
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
