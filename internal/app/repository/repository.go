package repository

import (
	"context"
	"fmt"
	"math/rand"

	appErrors "github.com/Te8va/shortURL/internal/app/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) (*URLRepository, error) {
	return &URLRepository{db: db}, nil
}

func (r *URLRepository) PingPg(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping: %w", err)
	}

	return nil
}

func (r *URLRepository) Save(ctx context.Context, url string) (string, error) {
	id := r.generateID()

	query := `WITH ins AS (
				INSERT INTO urlshrt (short, original) 
				VALUES ($1, $2)
				ON CONFLICT (original) DO NOTHING
				RETURNING short
			  )
			  SELECT short FROM ins
			  UNION ALL
			  SELECT short FROM urlshrt WHERE original = $2 LIMIT 1;`

	var existingShort string
	err := r.db.QueryRow(ctx, query, id, url).Scan(&existingShort)

	if err != nil {
		return "", fmt.Errorf("ошибка при сохранении или получении short URL: %w", err)
	}

	if existingShort != id {
		return existingShort, appErrors.ErrURLExists
	}

	return existingShort, nil
}

func (r *URLRepository) Get(ctx context.Context, id string) (string, bool) {
	query := `SELECT original FROM urlshrt WHERE short = $1;`

	var url string
	err := r.db.QueryRow(ctx, query, id).Scan(&url)
	if err == nil {
		return url, true
	}

	return "", false
}

func (r *URLRepository) SaveBatch(ctx context.Context, urls map[string]string) (map[string]string, error) {
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
		query := `INSERT INTO urlshrt (short, original) VALUES ($1, $2);`

		_, err := tx.Exec(ctx, query, id, originalURL)
		if err != nil {
			return nil, fmt.Errorf("ошибка сохранения URL в БД: %w", err)
		}

		result[correlationID] = id
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
