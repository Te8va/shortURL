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

func NewRouter(cfg *config.Config, saver service.URLSaver, getter service.URLGetter, pinger service.Pinger, deleter service.URLDelete) chi.Router {
	r := chi.NewRouter()

	if err := middleware.Initialize("info"); err != nil {
		log.Println("Failed to initialize middleware:", err)
	}

	r.Use(middleware.AuthMiddleware(cfg.JWTKey))
	r.Use(middleware.WithLogging)

	r.Mount("/api/urls", newURLRouter(saver, getter, cfg))
	r.Mount("/ping", newPingRouter(pinger))
	r.Mount("/api/user", newUserRouter(getter, deleter, cfg))

	return r
}

func newURLRouter(saver service.URLSaver, getter service.URLGetter, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	saveHandler := handler.NewSaveHandler(saver)
	getHandler := handler.NewGetterHandler(getter, cfg)

	r.Post("/", saveHandler.PostHandler)
	r.Get("/{id}", getHandler.GetHandler)
	r.Post("/shorten", saveHandler.PostHandlerJSON)
	r.Post("/batch", saveHandler.PostHandlerBatch)
	return r
}

func newPingRouter(pinger service.Pinger) chi.Router {
	r := chi.NewRouter()

	if pinger != nil {
		pingHandler := handler.NewPingHandler(pinger)
		r.Get("/", pingHandler.PingHandler)
	}

	return r
}

func newUserRouter(getter service.URLGetter, deleter service.URLDelete, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	getHandler := handler.NewGetterHandler(getter, cfg)
	deleteHandler := handler.NewDeleteHandler(deleter, cfg)

	r.Get("/urls", getHandler.GetUserURLsHandler)
	r.Delete("/urls", deleteHandler.DeleteUserURLsHandler)

	return r
}
