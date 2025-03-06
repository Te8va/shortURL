package router

import (
	"log"
	"github.com/go-chi/chi/v5"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/Te8va/shortURL/internal/app/repository"
)

func NewRouter(cfg *config.Config) chi.Router {
	repo, err := repository.NewMapStore(cfg.FileStoragePath)
	if err != nil {
		log.Println("Failed to initialize file repository:", err)
	}

	store := handler.NewURLStore(cfg, repo)
	r := chi.NewRouter()

	if err := middleware.Initialize("info"); err != nil {
		log.Println("Failed to initialize middleware:", err)
	}

	r.Use(middleware.WithLogging)
	r.Post("/", store.PostHandler)
	r.Get("/{id}", store.GetHandler)
	r.Post("/api/shorten", store.PostHandlerJSON)

	return r
}
