package handler

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const length = 8

type URLStore struct {
	store map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		store: make(map[string]string),
	}
}

func (u *URLStore) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		u.PostHandler(w, r)
	} else if r.Method == http.MethodGet {
		u.GetHandler(w, r)
	} else {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}

func (u *URLStore) PostHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	scanner := bufio.NewScanner(r.Body)
	if !scanner.Scan() {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalURL := scanner.Text()
	if originalURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	id := u.generateID()
	u.store[id] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(shortenedURL)); err != nil {
		http.Error(w, "Failed to write response", http.StatusBadRequest)
		return
	}
}

func (u *URLStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Missing or invalid ID in the URL path", http.StatusBadRequest)
		return
	}

	originalURL, exists := u.store[id]
	if !exists {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *URLStore) generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for {
		randStrBytes := make([]byte, length)
		for i := 0; i < length; i++ {
			randStrBytes[i] = charset[rand.Intn(len(charset))]
		}
		id := string(randStrBytes)

		if _, exists := u.store[id]; !exists {
			return id
		}
	}
}
