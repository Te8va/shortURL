package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const length = 8

type URLRepository struct {
	db   *pgxpool.Pool
	data map[string]string
	file string
}

func NewURLRepository(db *pgxpool.Pool, filePath string) (*URLRepository, error) {
	store := &URLRepository{
		db:   db,
		data: make(map[string]string),
		file: filePath,
	}

	if err := store.createDB(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.WriteFile(filePath, []byte("{}"), 0666)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания файла %s: %w", filePath, err)
		}
	}

	if err := store.loadFromFile(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки данных из файла %s: %w", filePath, err)
	}

	return store, nil
}

func (r *URLRepository) createDB(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS urlshrt (
		uuid SERIAL PRIMARY KEY,
		short TEXT UNIQUE NOT NULL,
		original TEXT NOT NULL
	)`)
	return err
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
	query := `INSERT INTO urlshrt (short, original) VALUES ($1, $2);`

	_, err := r.db.Exec(ctx, query, id, url)
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения в БД: %w", err)
	}

	r.data[id] = url

	err = r.saveToFile()
	if err != nil {
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

	url, exists := r.data[id]
	return url, exists
}

func (r *URLRepository) saveToFile() error {
	file, err := os.OpenFile(r.file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (r *URLRepository) loadFromFile() error {
	file, err := os.ReadFile(r.file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, &r.data); err != nil {
		return err
	}
	return nil
}

func (r *URLRepository) generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for {
		randStrBytes := make([]byte, length)
		for i := 0; i < length; i++ {
			randStrBytes[i] = charset[rand.Intn(len(charset))]
		}
		id := string(randStrBytes)

		if _, exists := r.data[id]; !exists {
			return id
		}
	}
}
