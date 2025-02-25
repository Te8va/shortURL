package main

import (
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/router"
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

	sugar.Infow(
		"Starting server",
		"addr", cfg.ServerAddress,
	)
	if err := http.ListenAndServe(cfg.ServerAddress, router.NewRouter(cfg)); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
