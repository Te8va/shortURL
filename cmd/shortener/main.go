package main

import (
	"log"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/router"
)

func main() {
	cfg := config.NewConfig()

	log.Println("Starting server on ", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router.NewRouter(cfg)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
