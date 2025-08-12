// package handler contains logic for retrieving original URLs and user-specific URL.

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

// URLGetter defines an interface for retrieving URLs.
//
//go:generate mockgen -source=gethandler.go -destination=mocks/url_getter_mock.gen.go -package=mocks
type URLGetter interface {
	Get(ctx context.Context, id string) (string, bool, bool)
	GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error)
}

// GetterHandler handles requests for retrieving URLs.
type GetterHandler struct {
	getter URLGetter
	cfg    *config.Config
}

// NewGetterHandler creates a new instance of GetterHandler.
func NewGetterHandler(getter URLGetter, cfg *config.Config) *GetterHandler {
	return &GetterHandler{getter: getter, cfg: cfg}
}

// GetHandler processes request to redirect to the original URL by short ID.
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

// GetUserURLsHandler a request to retrieve all URLs created user.
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

	w.Header().Set(contentType, contentTypeApp)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}
