package main

import (
	"log"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/Te8va/shortURL/internal/app/repository"
	"github.com/Te8va/shortURL/internal/app/service"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorw("Failed to sync logger", "error", err)
		}
	}()

	cfg := config.NewConfig()

	if err := env.Parse(cfg); err != nil {
		sugar.Fatalw("Failed to parse env", "error", err)
	}

	m, err := migrate.New("file://migrations", cfg.PostgresConn)
	if err != nil {
		sugar.Fatalw("Failed to apply migrations", "error", err)
	}

	err = repository.ApplyMigrations(m)
	if err != nil {
		sugar.Fatalw("Failed to apply migrations", "error", err)
	}

	sugar.Infoln("Migrations applied successfully")

	pool, err := repository.GetPgxPool(cfg.PostgresConn)
	if err != nil {
		sugar.Fatalln("Failed to connect to database", "error", err)
	}

	repo := repository.NewURLService(pool)
	srv := service.NewURL(repo)
	store := handler.NewURLStore(cfg, srv)
	r := chi.NewRouter()

	if err := middleware.Initialize("info"); err != nil {
		log.Println("Failed to initialize middleware:", err)
	}

	r.Use(middleware.WithLogging)
	r.Post("/", store.PostHandler)
	r.Get("/{id}", store.GetHandler)
	r.Post("/api/shorten", store.PostHandlerJSON)
	r.Get("/ping", store.PingHandler)

	sugar.Infow("Starting server", "addr", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
