// package handler contains handlers for saving new URLs, including single and batch operations.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Te8va/shortURL/internal/app/domain"
	appErrors "github.com/Te8va/shortURL/internal/app/errors"
)

const (
	length          = 8
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeApp  = "application/json"
)

// URLSaver defines an interface for saving URLs.
//
//go:generate mockgen -source=savehandler.go -destination=mocks/url_saver_mock.gen.go -package=mocks
type URLSaver interface {
	Save(ctx context.Context, userID int, url string) (string, error)
	SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error)
}

// SaveHandler handles requests for saving URLs.
type SaveHandler struct {
	saver URLSaver
}

// NewSaveHandler creates a new instance of SaveHandler.
func NewSaveHandler(saver URLSaver) *SaveHandler {
	return &SaveHandler{saver: saver}
}

// PostHandler processes requests to save URL.
func (u *SaveHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get(ContentType), ContentTypeText) {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(domain.UserIDKey).(int)

	originalURLBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	originalURL := string(originalURLBytes)
	if originalURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := u.saver.Save(ctx, userID, originalURL)
	if err != nil {
		if !errors.Is(err, appErrors.ErrURLExists) {
			http.Error(w, "Failed to save URL", http.StatusBadRequest)
			return
		}
		w.Header().Set(ContentType, ContentTypeText)
		w.WriteHeader(http.StatusConflict)
		if _, err := w.Write([]byte(id)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set(ContentType, ContentTypeText)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(id)); err != nil {
		http.Error(w, "Failed to write response", http.StatusBadRequest)
		return
	}
}

// PostHandlerJSON processes JSON-formatted POST requests to save URL.
func (u *SaveHandler) PostHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get(ContentType), ContentTypeApp) {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(domain.UserIDKey).(int)

	var req domain.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	id, err := u.saver.Save(r.Context(), userID, req.URL)

	if errors.Is(err, appErrors.ErrURLExists) {
		resp := domain.ShortenResponse{Result: id}
		w.Header().Set(ContentType, ContentTypeApp)
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	} else if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	resp := domain.ShortenResponse{Result: id}

	w.Header().Set(ContentType, ContentTypeApp)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// BatchRequest represents an individual request item in a batch of URLs to shorten
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse represents a shortened URL response for a single batch item.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// PostHandlerBatch processes batch URL saving requests in JSON format.
func (u *SaveHandler) PostHandlerBatch(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(domain.UserIDKey).(int)

	var batchReq []BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&batchReq); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if len(batchReq) == 0 {
		http.Error(w, "Пустой список URL", http.StatusBadRequest)
		return
	}

	urlMap := make(map[string]string)
	for _, req := range batchReq {
		id, err := u.saver.Save(r.Context(), userID, req.OriginalURL)
		if err != nil {
			http.Error(w, "Ошибка сохранения URL", http.StatusInternalServerError)
			return
		}
		urlMap[req.CorrelationID] = id
	}

	var batchResp []BatchResponse
	for _, req := range batchReq {
		batchResp = append(batchResp, BatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      urlMap[req.CorrelationID],
		})
	}

	w.Header().Set(ContentType, ContentTypeApp)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(batchResp); err != nil {
		http.Error(w, "Ошибка записи ответа", http.StatusInternalServerError)
	}
}
