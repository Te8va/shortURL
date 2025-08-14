package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/Te8va/shortURL/internal/app/handler"
)

type mockDeleter struct{}

func (m mockDeleter) DeleteUserURLs(ctx context.Context, ids []string, userID int) error {
	return nil
}

func ExampleDeleteHandler_DeleteUserURLsHandler() {
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	h := handler.NewDeleteHandler(mockDeleter{}, cfg)

	r := chi.NewRouter()
	r.Delete("/user/urls", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), domain.UserIDKey, 1)
		h.DeleteUserURLsHandler(w, r.WithContext(ctx))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	ids := []string{"abc123", "xyz789"}
	body, _ := json.Marshal(ids)

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/user/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 202
}
