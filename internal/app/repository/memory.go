package repository

import (
	"context"
	"math/rand"
)

type MemoryRepository struct {
	store map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		store: make(map[string]string),
	}
}

func (r *MemoryRepository) Save(ctx context.Context, userID int, url string) (string, error) {
	id := r.generateID()
	r.store[id] = url

	return id, nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (string, bool) {

	url, exists := r.store[id]
	return url, exists
}

func (r *MemoryRepository) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	for correlationID, originalURL := range urls {
		id := r.generateID()
		r.store[id] = originalURL
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

func (r *MemoryRepository) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	return []map[string]string{}, nil
}
