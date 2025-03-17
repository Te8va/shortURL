package router

import (
	"log"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/Te8va/shortURL/internal/app/service"
)

func NewRouter(cfg *config.Config, saver service.URLSaver, getter service.URLGetter, pinger service.Pinger) chi.Router {
	srv := service.NewURLService(saver, getter, pinger)
	store := handler.NewURLHandler(cfg, srv, srv, srv)
	r := chi.NewRouter()

	if err := middleware.Initialize("info"); err != nil {
		log.Println("Failed to initialize middleware:", err)
	}

	r.Use(middleware.WithLogging)
	r.Post("/", store.PostHandler)
	r.Get("/{id}", store.GetHandler)
	r.Post("/api/shorten", store.PostHandlerJSON)
	r.Post("/api/shorten/batch", store.PostHandlerBatch)
	r.Get("/ping", store.PingHandler)

	return r
}
