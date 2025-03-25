package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"

	"github.com/Te8va/shortURL/internal/app/config"
	appErrors "github.com/Te8va/shortURL/internal/app/errors"
)

const length = 8

type JSONRepository struct {
	file  string
	store map[string]URLData
	mu    sync.RWMutex
	cfg   *config.Config
}

type URLData struct {
	UserID      int    `json:"user_id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func NewJSONRepository(filePath string, cfg *config.Config) (*JSONRepository, error) {
	if filePath == "" {
		return nil, fmt.Errorf("путь к файлу не задан")
	}

	repo := &JSONRepository{
		file:  filePath,
		store: make(map[string]URLData),
		cfg:   cfg,
	}

	if err := repo.loadFromFile(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки данных из файла: %w", err)
	}

	return repo, nil
}

func (r *JSONRepository) Save(ctx context.Context, userID int, url string) (string, error) {
	id := r.generateID()
	shortenedURL := fmt.Sprintf("%s/%s", r.cfg.BaseURL, id)

	r.mu.Lock()
	defer r.mu.Unlock()

	for key, val := range r.store {
		if val.OriginalURL == url && val.UserID == userID {
			return key, appErrors.ErrURLExists
		}
	}

	r.store[shortenedURL] = URLData{
		UserID:      userID,
		OriginalURL: url,
		ShortURL:    shortenedURL,
	}

	if err := r.saveToFile(); err != nil {
		return "", fmt.Errorf("ошибка сохранения в файл: %w", err)
	}

	return shortenedURL, nil
}

func (r *JSONRepository) Get(ctx context.Context, id string, errChan chan error) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.store[id]
	if !exists {
		return "", appErrors.ErrNotFound
	}
	return url.OriginalURL, nil
}

func (r *JSONRepository) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	r.mu.Lock()
	defer r.mu.Unlock()

	for correlationID, originalURL := range urls {
		id := r.generateID()
		shortenedURL := fmt.Sprintf("%s/%s", r.cfg.BaseURL, id)
		r.store[shortenedURL] = URLData{
			UserID:      userID,
			OriginalURL: originalURL,
			ShortURL:    shortenedURL,
		}
		result[correlationID] = shortenedURL
	}

	if err := r.saveToFile(); err != nil {
		return nil, fmt.Errorf("ошибка сохранения в файл: %w", err)
	}

	return result, nil
}

func (r *JSONRepository) saveToFile() error {
	file, err := os.OpenFile(r.file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла %s: %w", r.file, err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(r.store, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %w", r.file, err)
	}

	return nil
}

func (r *JSONRepository) loadFromFile() error {
	fileData, err := os.ReadFile(r.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if len(fileData) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := json.Unmarshal(fileData, &r.store); err != nil {
		return fmt.Errorf("ошибка десериализации данных из файла: %w", err)
	}

	return nil
}

func (r *JSONRepository) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var urls []map[string]string
	for _, data := range r.store {
		if data.UserID == userID {
			urls = append(urls, map[string]string{
				"short_url":    data.ShortURL,
				"original_url": data.OriginalURL,
			})
		}
	}

	if len(urls) == 0 {
		return nil, nil
	}

	return urls, nil
}

func (r *JSONRepository) generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for {
		randStrBytes := make([]byte, length)
		for i := 0; i < length; i++ {
			randStrBytes[i] = charset[rand.Intn(len(charset))]
		}
		id := string(randStrBytes)

		r.mu.RLock()
		_, exists := r.store[id]
		r.mu.RUnlock()

		if !exists {
			return id
		}
	}
}
