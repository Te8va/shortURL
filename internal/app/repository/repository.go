package repository

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/Te8va/shortURL/internal/app/config"
	appErrors "github.com/Te8va/shortURL/internal/app/errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewURLRepository(db *pgxpool.Pool, cfg *config.Config) (*URLRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	return &URLRepository{db: db, cfg: cfg}, nil
}

func (r *URLRepository) PingPg(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping: %w", err)
	}

	return nil
}

func (r *URLRepository) Save(ctx context.Context, userID int, url string) (string, error) {
	id := r.generateID()
	shortenedURL := fmt.Sprintf("%s/%s", r.cfg.BaseURL, id)

	query := `WITH ins AS (
				INSERT INTO urlshrt (short, original, user_id) 
				VALUES ($1, $2, $3)
				ON CONFLICT (original) DO NOTHING
				RETURNING short
			  )
			  SELECT short FROM ins
			  UNION ALL
			  SELECT short FROM urlshrt WHERE original = $2 LIMIT 1;`

	var existingShort string
	err := r.db.QueryRow(ctx, query, shortenedURL, url, userID).Scan(&existingShort)

	if err != nil {
		return "", fmt.Errorf("ошибка при сохранении или получении short URL: %w", err)
	}

	if existingShort != shortenedURL {
		return existingShort, appErrors.ErrURLExists
	}

	return existingShort, nil
}

func (r *URLRepository) Get(ctx context.Context, id string) (string, error) {
	query := `SELECT original, is_deleted FROM urlshrt WHERE short = $1;`

	var url string
	var isDeleted bool

	err := r.db.QueryRow(ctx, query, id).Scan(&url, &isDeleted)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", appErrors.ErrNotFound
		}
		return "", fmt.Errorf("ошибка запроса в БД: %w", err)
	}

	if isDeleted {
		return "", appErrors.ErrDeleted
	}

	return url, nil
}

func (r *URLRepository) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && err == nil {
			err = fmt.Errorf("ошибка при откате транзакции: %w", rollbackErr)
		}
	}()

	result := make(map[string]string)
	for correlationID, originalURL := range urls {
		id := r.generateID()
		shortenedURL := fmt.Sprintf("%s/%s", r.cfg.BaseURL, id)
		query := `INSERT INTO urlshrt (short, original, user_id) VALUES ($1, $2, $3);`

		_, err := tx.Exec(ctx, query, shortenedURL, originalURL)
		if err != nil {
			return nil, fmt.Errorf("ошибка сохранения URL в БД: %w", err)
		}

		result[correlationID] = shortenedURL
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("ошибка при завершении транзакции: %w", err)
	}

	return result, nil
}

func (r *URLRepository) generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for {
		randStrBytes := make([]byte, length)
		for i := 0; i < length; i++ {
			randStrBytes[i] = charset[rand.Intn(len(charset))]
		}
		id := string(randStrBytes)

		var exists bool
		err := r.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM urlshrt WHERE short = $1)", id).Scan(&exists)
		if err != nil {
			continue
		}

		if !exists {
			return id
		}
	}
}

func (r *URLRepository) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	query := `SELECT short, original FROM urlshrt WHERE user_id = $1;`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении URL пользователя: %w", err)
	}
	defer rows.Close()

	var urls []map[string]string
	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании URL: %w", err)
		}
		urls = append(urls, map[string]string{
			"short_url":    shortURL,
			"original_url": originalURL,
		})
	}

	if len(urls) == 0 {
		return nil, nil
	}

	return urls, nil
}

func (r *URLRepository) DeleteUserURLs(ctx context.Context, userID int, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	const workerCount = 4
	const batchSize = 50

	inputCh := make(chan []string)
	errCh := make(chan error)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range inputCh {
				if err := r.DeleteUserURL(ctx, userID, batch); err != nil {
					errCh <- err
				}
			}
		}()
	}

	go func() {
		defer close(inputCh)
		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			inputCh <- ids[i:end]
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var finalErr error
	for err := range errCh {
		if finalErr == nil {
			finalErr = err
		}
	}

	return finalErr
}

func (r *URLRepository) DeleteUserURL(ctx context.Context, userID int, ids []string) error {
	query := `UPDATE urlshrt SET is_deleted = true WHERE user_id = $1 AND short = ANY($2);`

	_, err := r.db.Exec(ctx, query, userID, ids)
	if err != nil {
		return fmt.Errorf("ошибка при удалении URL: %w", err)
	}
	return nil
}
