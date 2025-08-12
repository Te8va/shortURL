package repository

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Te8va/shortURL/internal/app/config"
)

// MemoryRepository is a storage implementation that keeps data in memory.
type MemoryRepository struct {
	store map[string]string
	cfg   *config.Config
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository(cfg *config.Config) *MemoryRepository {
	return &MemoryRepository{
		store: make(map[string]string),
		cfg:   cfg,
	}
}

// Save stores URL and returns its shortened version.
func (r *MemoryRepository) Save(ctx context.Context, userID int, url string) (string, error) {
	id := r.generateID()
	shortenedURL := fmt.Sprintf("%s/%s", r.cfg.BaseURL, id)
	r.store[shortenedURL] = url

	return shortenedURL, nil
}

// Get returns the original URL by its shortened identifier
func (r *MemoryRepository) Get(ctx context.Context, id string) (string, bool, bool) {
	url, exists := r.store[id]
	if !exists {
		return "", false, false
	}
	return url, true, false
}

// SaveBatch stores multiple URLs in a single call
func (r *MemoryRepository) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	for correlationID, originalURL := range urls {
		id := r.generateID()
		shortenedURL := fmt.Sprintf("%s/%s", r.cfg.BaseURL, id)

		r.store[shortenedURL] = originalURL
		result[correlationID] = id
	}

	return result, nil
}

func (r *MemoryRepository) generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for {
		randStrBytes := make([]byte, length)
		for i := 0; i < length; i++ {
			randStrBytes[i] = charset[rand.Intn(len(charset))]
		}
		id := string(randStrBytes)

		_, exists := r.store[id]

		if !exists {
			return id
		}
	}
}

// GetUserURLs returns all URLs belonging to a specific user
func (r *MemoryRepository) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	return []map[string]string{}, nil
}
