package repository

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

const length = 8

type URLService struct {
	pool *pgxpool.Pool
}

func NewURLService(pool *pgxpool.Pool) *URLService {
	return &URLService{pool: pool}
}

func (r *URLService) PingPg(ctx context.Context) error {
	err := r.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping: %w", err)
	}
	return nil
}

func (r *URLService) Save(ctx context.Context, url string) (string, error) {
	id := r.generateID()

	_, err := r.pool.Exec(ctx, "INSERT INTO urls (id, url) VALUES ($1, $2)", id, url)
	if err != nil {
		return "", fmt.Errorf("ошибка при сохранении URL: %w", err)
	}

	return id, nil
}

func (r *URLService) Get(ctx context.Context, id string) (string, bool) {
	var url string
	err := r.pool.QueryRow(ctx, "SELECT url FROM urls WHERE id = $1", id).Scan(&url)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", false
		}
		return "", false
	}

	return url, true
}

func (r *URLService) generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for {
		randStrBytes := make([]byte, length)
		for i := 0; i < length; i++ {
			randStrBytes[i] = charset[rand.Intn(len(charset))]
		}
		id := string(randStrBytes)

		var exists bool
		err := r.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM urls WHERE id = $1)", id).Scan(&exists)
		if err != nil {
			continue
		}

		if !exists {
			return id
		}
	}
}
