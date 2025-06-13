package router

import (
	"log"

	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	mdlwr "github.com/go-chi/chi/v5/middleware"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/Te8va/shortURL/internal/app/service"
)

func NewRouter(cfg *config.Config, saver service.URLSaverServ, getter service.URLGetterServ, pinger service.PingerServ, deleter service.URLDeleteServ) chi.Router {
	r := chi.NewRouter()

	if err := middleware.Initialize("info"); err != nil {
		log.Println("Failed to initialize middleware:", err)
	}

	r.Use(middleware.AuthMiddleware(cfg.JWTKey))
	r.Use(middleware.WithLogging)

	r.Mount("/", newRootRouter(cfg, saver, getter))
	r.Mount("/api", newAPIRouter(cfg, saver, getter, deleter))
	r.Mount("/ping", newPingRouter(pinger))
	r.Mount("/debug", mdlwr.Profiler())

	return r
}

func newRootRouter(cfg *config.Config, saver service.URLSaverServ, getter service.URLGetterServ) chi.Router {
	r := chi.NewRouter()

	saveHandler := handler.NewSaveHandler(saver)
	getHandler := handler.NewGetterHandler(getter, cfg)

	r.Post("/", saveHandler.PostHandler)
	r.Get("/{id}", getHandler.GetHandler)

	return r
}

func newAPIRouter(cfg *config.Config, saver service.URLSaverServ, getter service.URLGetterServ, deleter service.URLDeleteServ) chi.Router {
	r := chi.NewRouter()

	saveHandler := handler.NewSaveHandler(saver)
	getHandler := handler.NewGetterHandler(getter, cfg)
	deleteHandler := handler.NewDeleteHandler(deleter, cfg)

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

func newPingRouter(pinger service.PingerServ) chi.Router {
	r := chi.NewRouter()

	if pinger != nil {
		pingHandler := handler.NewPingHandler(pinger)
		r.Get("/", pingHandler.PingHandler)
	}

	return r
}
