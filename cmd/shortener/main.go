package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
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

	var storage domain.RepositoryStore

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.DatabaseDSN != "" {
		sugar.Infoln("Using PostgreSQL as storage")

		m, err := migrate.New("file://migrations", cfg.DatabaseDSN)
		if err != nil {
			sugar.Fatalw("Failed to initialize migrations", "error", err)
		}

		err = repository.ApplyMigrations(m)
		if err != nil {
			sugar.Fatalw("Failed to apply migrations", "error", err)
		}

		pool, err := repository.GetPgxPool(ctx, cfg.DatabaseDSN)
		if err != nil {
			sugar.Fatalw("Failed to create Postgres connection pool", "error", err)
		}
		defer pool.Close()

		storage, err = repository.NewURLRepository(pool)
		if err != nil {
			sugar.Fatalw("Failed to initialize Postgres repository", "error", err)
		}

	} else if cfg.FileStoragePath != "" {
		sugar.Infoln("Using JSON file as storage:", cfg.FileStoragePath)

		storage, err = repository.NewJSONRepository(cfg.FileStoragePath)
		if err != nil {
			sugar.Fatalw("Failed to initialize JSON repository", "error", err)
		}

	} else {
		sugar.Infoln("Using in-memory storage")
		storage = repository.NewMemoryRepository()
	}

	handler := router.NewRouter(cfg, storage)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}

	go func() {
		sugar.Infow("Server started", "addr", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("ListenAndServe failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	sugar.Infoln("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		sugar.Fatalw("Server shutdown failed", "error", err)
	}

	var wg sync.WaitGroup
	waitGroupChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitGroupChan)
	}()

	select {
	case <-waitGroupChan:
		sugar.Infoln("All goroutines finished cleanly")
	case <-time.After(3 * time.Second):
		cancel()
		sugar.Infoln("Some goroutines did not finish in time")
	}

	sugar.Infoln("Server shut down successfully")
}
