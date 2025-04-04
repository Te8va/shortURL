package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
)

type URLDelete interface {
	DeleteUserURLs(ctx context.Context, ids []string, userID int) error
}

type DeleteHandler struct {
	deleter URLDelete
	cfg     *config.Config
}

func NewDeleteHandler(deleter URLDelete, cfg *config.Config) *DeleteHandler {
	return &DeleteHandler{deleter: deleter, cfg: cfg}
}

func (u *DeleteHandler) DeleteUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(ids) == 0 {
		http.Error(w, "Empty list of URLs", http.StatusBadRequest)
		return
	}

	var fullURLs []string
	for _, id := range ids {
		fullURLs = append(fullURLs, fmt.Sprintf("%s/%s", u.cfg.BaseURL, id))
	}

	go func(fullURLs []string, userID int) {
		err := u.deleter.DeleteUserURLs(context.Background(), fullURLs, userID)
		if err != nil {
			http.Error(w, "Failed to delete URL", http.StatusBadRequest)
			return
		}
	}(fullURLs, userID)

	w.WriteHeader(http.StatusAccepted)
}
