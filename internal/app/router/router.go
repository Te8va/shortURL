package router

// import (
// 	"log"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	_ "github.com/jackc/pgx/v5/stdlib"

// 	"github.com/Te8va/shortURL/internal/app/config"
// 	"github.com/Te8va/shortURL/internal/app/handler"
// 	"github.com/Te8va/shortURL/internal/app/middleware"
// 	"github.com/Te8va/shortURL/internal/app/repository"
// 	"github.com/Te8va/shortURL/internal/app/service"
// )

// func NewRouter(cfg *config.Config, db *pgxpool.Pool) chi.Router {
// 	repo, err := repository.NewURLRepository(db, cfg.FileStoragePath)
// 	if err != nil {
// 		log.Println("Failed to initialize file repository:", err)
// 	}
// 	srv := service.NewURLService(repo)
// 	store := handler.NewURLStore(cfg, srv)
// 	r := chi.NewRouter()

// 	if err := middleware.Initialize("info"); err != nil {
// 		log.Println("Failed to initialize middleware:", err)
// 	}

// 	r.Use(middleware.WithLogging)
// 	r.Post("/", store.PostHandler)
// 	r.Get("/{id}", store.GetHandler)
// 	r.Post("/api/shorten", store.PostHandlerJSON)
// 	r.Get("/ping", store.PingHandler)

// 	return r
// }
