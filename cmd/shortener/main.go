package main

import (
	"context"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/router"
	"github.com/jackc/pgx/v5/pgxpool"
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

	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.PostgresConn)
	if err != nil {
		sugar.Errorw("Failed to connect to database: %v", "error", err)
	}
	defer db.Close()

	sugar.Infow(
		"Starting server",
		"addr", cfg.ServerAddress,
	)
	if err := http.ListenAndServe(cfg.ServerAddress, router.NewRouter(cfg, db)); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
