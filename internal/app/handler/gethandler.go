package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
)

type URLGetter interface {
	Get(ctx context.Context, id string) (string, bool, bool)
	GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error)
}

type GetterHandler struct {
	getter URLGetter
	cfg    *config.Config
}

func NewGetterHandler(getter URLGetter, cfg *config.Config) *GetterHandler {
	return &GetterHandler{getter: getter, cfg: cfg}
}

func (u *GetterHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Missing or invalid ID", http.StatusBadRequest)
		return
	}

	id = fmt.Sprintf("%s/%s", u.cfg.BaseURL, id)
	originalURL, exists, isDeleted := u.getter.Get(r.Context(), id)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	if isDeleted {
		http.Error(w, "URL has been deleted", http.StatusGone)
		return
	}

	log.Printf("Redirecting ID %s to URL: %s", id, originalURL)
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *GetterHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := u.getter.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if urls == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set(ContentType, ContentTypeApp)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}
