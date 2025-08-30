package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/repository"
	"github.com/Te8va/shortURL/internal/app/router"
	"github.com/Te8va/shortURL/internal/app/service"
)

// App represents the core application structure
type App struct {
	cfg     *config.Config
	logger  *zap.SugaredLogger
	saver   service.URLSaverServ
	getter  service.URLGetterServ
	pinger  service.PingerServ
	deleter service.URLDeleteServ
	server  *http.Server
}

// NewApp creates a new App instance
func NewApp() (*App, error) {
	cfg := config.NewConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorw("Failed to sync logger", "error", err)
		}
	}()

	app := &App{
		cfg:    cfg,
		logger: sugar,
	}

	if err := app.initStorage(); err != nil {
		return nil, err
	}

	app.initServer()
	return app, nil
}

func (a *App) initStorage() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch {
	case a.cfg.DatabaseDSN != "":
		return a.initPostgresStorage(ctx)
	case a.cfg.FileStoragePath != "":
		return a.initFileStorage()
	default:
		return a.initMemoryStorage()
	}
}

func (a *App) initPostgresStorage(ctx context.Context) error {
	a.logger.Infoln("Using PostgreSQL as storage")

	m, err := migrate.New("file://migrations", a.cfg.DatabaseDSN)
	if err != nil {
		a.logger.Fatalw("Failed to initialize migrations", "error", err)
	}

	if err = repository.ApplyMigrations(m); err != nil {
		a.logger.Fatalw("Failed to apply migrations", "error", err)
	}

	pool, err := repository.GetPgxPool(ctx, a.cfg.DatabaseDSN)
	if err != nil {
		a.logger.Fatalw("Failed to create Postgres connection pool", "error", err)
	}

	repo, err := repository.NewURLRepository(pool, a.cfg)
	if err != nil {
		a.logger.Fatalw("Failed to initialize Postgres repository", "error", err)
	}

	a.saver = repo
	a.getter = repo
	a.pinger = repo
	a.deleter = repo

	return nil
}

func (a *App) initFileStorage() error {
	a.logger.Infoln("Using JSON file as storage:", a.cfg.FileStoragePath)

	storage, err := repository.NewJSONRepository(a.cfg.FileStoragePath, a.cfg)
	if err != nil {
		a.logger.Fatalw("Failed to initialize JSON repository", "error", err)
	}

	a.saver = storage
	a.getter = storage
	return nil
}

func (a *App) initMemoryStorage() error {
	a.logger.Infoln("Using in-memory storage")
	storage := repository.NewMemoryRepository(a.cfg)

	a.saver = storage
	a.getter = storage
	return nil
}

func (a *App) initServer() {
	handler := router.NewRouter(a.cfg, a.saver, a.getter, a.pinger, a.deleter)

	a.server = &http.Server{
		Addr:    a.cfg.ServerAddress,
		Handler: handler,
	}
}

// Run starts the HTTP server and listens for OS signals to gracefully shut down all resources before exiting
func (a *App) Run() error {
	defer func() {
		if err := a.logger.Sync(); err != nil {
			a.logger.Errorw("Failed to sync logger", "error", err)
		}
	}()

	go func() {
		a.logger.Infow("Server started", "addr", a.cfg.ServerAddress)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatalw("ListenAndServe failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	a.logger.Infoln("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Fatalw("Server shutdown failed", "error", err)
	}

	var wg sync.WaitGroup
	waitGroupChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitGroupChan)
	}()

	select {
	case <-waitGroupChan:
		a.logger.Infoln("All goroutines finished cleanly")
	case <-time.After(3 * time.Second):
		a.logger.Warn("Some goroutines did not finish in time")
	}

	a.logger.Infoln("Server shut down successfully")
	return nil
}
