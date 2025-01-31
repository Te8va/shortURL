package router

import (
	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	store := handler.NewURLStore()
	r := chi.NewRouter()

	r.Use(middleware.Middleware)
	r.Post("/", store.PostHandler)
	r.Get("/{id}", store.GetHandler)

	return r
}