package main

import (
	"log"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/router"
)

func main() {
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router.NewRouter()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
