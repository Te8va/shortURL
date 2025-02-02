package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
)

func NewRouter(cfg *config.Config) chi.Router {
	store := handler.NewURLStore(cfg)
	r := chi.NewRouter()

	r.Use(middleware.Middleware)
	r.Post("/", store.PostHandler)
	r.Get("/{id}", store.GetHandler)

	return r
}
