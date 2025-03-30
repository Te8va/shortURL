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

	r.Mount("/", newRootRouter(saver, getter))
	r.Mount("/api", newAPIRouter(saver, getter, deleter))
	r.Mount("/ping", newPingRouter(pinger))

	return r
}

func newRootRouter(saver service.URLSaver, getter service.URLGetter) chi.Router {
	r := chi.NewRouter()

	saveHandler := handler.NewSaveHandler(saver)
	getHandler := handler.NewGetterHandler(getter, nil)

	r.Post("/", saveHandler.PostHandler)
	r.Get("/{id}", getHandler.GetHandler)

	return r
}

func newAPIRouter(saver service.URLSaver, getter service.URLGetter, deleter service.URLDelete) chi.Router {
	r := chi.NewRouter()

	saveHandler := handler.NewSaveHandler(saver)
	getHandler := handler.NewGetterHandler(getter, nil)
	deleteHandler := handler.NewDeleteHandler(deleter, nil)

	r.Route("/shorten", func(r chi.Router) {
		r.Post("/", saveHandler.PostHandlerJSON)
		r.Post("/batch", saveHandler.PostHandlerBatch)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/urls", getHandler.GetUserURLsHandler)
		r.Delete("/urls", deleteHandler.DeleteUserURLsHandler)
	})

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
