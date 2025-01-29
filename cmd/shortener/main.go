package main

import (
	"log"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/handler"
	"github.com/Te8va/shortURL/internal/app/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Middleware(http.HandlerFunc(handler.NewURLStore().RootHandler)))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
