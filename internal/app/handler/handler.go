package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
	appErrors "github.com/Te8va/shortURL/internal/app/errors"
)

const (
	length          = 8
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeApp  = "application/json"
)

type URLSaver interface {
	Save(ctx context.Context, userID int, url string) (string, error)
	SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error)
}

type URLGetter interface {
	Get(ctx context.Context, id string, errChan chan error) (string, error)
	GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error)
}

type URLDelete interface {
	DeleteUserURLs(ctx context.Context, userID int, ids []string) error
}

type Pinger interface {
	PingPg(ctx context.Context) error
}

type URLHandler struct {
	saver   URLSaver
	getter  URLGetter
	pinger  Pinger
	deleter URLDelete
	cfg     *config.Config
}

func NewURLHandler(cfg *config.Config, saver URLSaver, getter URLGetter, pinger Pinger, deleter URLDelete) *URLHandler {
	return &URLHandler{
		saver:   saver,
		getter:  getter,
		pinger:  pinger,
		deleter: deleter,
		cfg:     cfg,
	}
}

func (u *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := u.pinger.PingPg(r.Context())
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (u *URLHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
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

func (u *URLHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Missing or invalid ID", http.StatusBadRequest)
		return
	}

	id = fmt.Sprintf("%s/%s", u.cfg.BaseURL, id)

	errChan := make(chan error, 1)
	originalURL, err := u.getter.Get(r.Context(), id, errChan)
	select {
	case errDeleted := <-errChan:
		log.Printf("URL %s был удален: %v", id, errDeleted)
		http.Error(w, "URL has been deleted", http.StatusGone)
		return
	default:
		if err != nil {
			if errors.Is(err, appErrors.ErrDeleted) {
				http.Error(w, "URL has been deleted", http.StatusGone)
				return
			}
			if errors.Is(err, appErrors.ErrNotFound) {
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Redirecting ID %s to URL: %s", id, originalURL)
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *URLHandler) PostHandlerJSON(w http.ResponseWriter, r *http.Request) {
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

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (u *URLHandler) PostHandlerBatch(w http.ResponseWriter, r *http.Request) {
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

func (u *URLHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

func (u *URLHandler) DeleteUserURLsHandler(w http.ResponseWriter, r *http.Request) {
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

	go func(ids []string, userID int) {
		err := u.deleter.DeleteUserURLs(context.Background(), userID, ids)
		if err != nil {
			log.Printf("Ошибка при удалении URL: %v", err)
		}
	}(ids, userID)

	w.WriteHeader(http.StatusAccepted)
}
