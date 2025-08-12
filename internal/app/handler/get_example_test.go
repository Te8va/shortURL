package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/Te8va/shortURL/internal/app/handler"
)

type mockGetter struct{}

func (m mockGetter) Get(ctx context.Context, id string) (string, bool, bool) {
	if strings.HasSuffix(id, "/abc123") {
		return "https://example.com", true, false
	}
	return "", false, false
}

func (m mockGetter) GetUserURLs(ctx context.Context, userID int) ([]map[string]string, error) {
	return []map[string]string{
		{"short_url": fmt.Sprintf("%s/%s", "http://example.test", "abc123"), "original_url": "https://example.com"},
	}, nil
}

func ExampleGetterHandler_GetHandler() {
	r := chi.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	cfg := &config.Config{BaseURL: ts.URL}
	h := handler.NewGetterHandler(mockGetter{}, cfg)

	r.Get("/{id}", h.GetHandler)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(ts.URL + "/abc123")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Location"))
	// Output:
	// 307
	// https://example.com
}

func ExampleGetterHandler_GetUserURLsHandler() {
	r := chi.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	cfg := &config.Config{BaseURL: ts.URL}
	h := handler.NewGetterHandler(mockGetter{}, cfg)

	r.Get("/user/urls", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), domain.UserIDKey, 1)
		h.GetUserURLsHandler(w, r.WithContext(ctx))
	})

	resp, err := http.Get(ts.URL + "/user/urls")
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	defer resp.Body.Close()

	var urls []map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&urls)

	fmt.Println(resp.StatusCode)
	fmt.Println(urls[0]["short_url"])
	// Output:
	// 200
	// http://example.test/abc123
}
