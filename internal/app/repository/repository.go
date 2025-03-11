package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"

	appErrors "github.com/Te8va/shortURL/internal/app/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgxpool"
)

const length = 8

type URLRepository struct {
	db   *pgxpool.Pool
	file string
}

func NewURLRepository(db *pgxpool.Pool, filePath string) (*URLRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("пул подключений к базе данных равен nil")
	}

	store := &URLRepository{
		db:   db,
		file: filePath,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.WriteFile(filePath, []byte("{}"), 0666)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания файла %s: %w", filePath, err)
		}
	}

	return store, nil
}

func (r *URLRepository) PingPg(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping: %w", err)
	}

	return nil
}

func (r *URLRepository) Save(ctx context.Context, url string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("URLRepository не инициализирован")
	}

	id := r.generateID()
	query := `INSERT INTO urlshrt (short, original) 
        VALUES ($1, $2) 
        ON CONFLICT (original) 
        DO UPDATE SET short = EXCLUDED.short 
        RETURNING short;`

	var short string
	err := r.db.QueryRow(ctx, query, id, url).Scan(&short)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			uErr := appErrors.ErrURLExists

			row := r.db.QueryRow(ctx, "SELECT short FROM urlshrt WHERE original = $1", url)
			var existingShort string
			errScan := row.Scan(&existingShort)
			if errScan != nil {
				return "", errScan
			}

			return existingShort, uErr
		}

		return "", fmt.Errorf("ошибка сохранения в БД: %w", err)
	}

	if short != id {
		return short, appErrors.ErrURLExists
	}

	if err := r.saveToFile(id, url); err != nil {
		return "", fmt.Errorf("ошибка сохранения в файл: %w", err)
	}

	return id, err
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

func (r *URLRepository) saveToFile(id, url string) error {
	file, err := os.OpenFile(r.file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла %s: %w", r.file, err)
	}
	defer file.Close()

	data := make(map[string]string)
	fileData, err := os.ReadFile(r.file)
	if err == nil && len(fileData) > 0 {
		if err := json.Unmarshal(fileData, &data); err != nil {
			return fmt.Errorf("ошибка десериализации данных из файла: %w", err)
		}
	}

	data[id] = url

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	_, err = file.WriteAt(jsonData, 0)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %w", r.file, err)
	}

	return err
}

func (r *URLRepository) SaveBatch(ctx context.Context, urls map[string]string) (map[string]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	result := make(map[string]string)
	for correlationID, originalURL := range urls {
		id := r.generateID()
		query := `INSERT INTO urlshrt (short, original) VALUES ($1, $2);`

		_, err := tx.Exec(ctx, query, id, originalURL)
		if err != nil {
			return nil, fmt.Errorf("ошибка сохранения URL в БД: %w", err)
		}

		if err := r.saveToFile(id, originalURL); err != nil {
			return nil, fmt.Errorf("ошибка сохранения в файл: %w", err)
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
