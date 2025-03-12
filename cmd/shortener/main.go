package main

import (
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/repository"
	"github.com/Te8va/shortURL/internal/app/router"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

func main() {

	cfg := config.NewConfig()

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

	sugar.Infoln(*cfg)

	m, err := migrate.New("file:///migrations", cfg.DatabaseDSN)
	if err != nil {
		sugar.Fatalw("Failed to initialize migrations", "error", err)
	}

	err = repository.ApplyMigrations(m)
	if err != nil {
		sugar.Fatalw("Failed to apply migrations", "error", err)
	}

	pool, err := repository.GetPgxPool(cfg.DatabaseDSN)
	if err != nil {
		sugar.Fatalw("Failed to create Postgres connection pool", "error", err)
	}

	defer pool.Close()

	sugar.Infow("Migrations applied successfully")

	sugar.Infow(
		"Starting server",
		"addr", cfg.ServerAddress,
	)
	if err := http.ListenAndServe(cfg.ServerAddress, router.NewRouter(cfg, pool)); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
