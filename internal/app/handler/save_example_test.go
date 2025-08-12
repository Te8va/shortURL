package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/Te8va/shortURL/internal/app/handler"
)

type mockSaver struct{}

func (m mockSaver) Save(ctx context.Context, userID int, url string) (string, error) {
	return "http://short.ly/abc123", nil
}
func (m mockSaver) SaveBatch(ctx context.Context, userID int, urls map[string]string) (map[string]string, error) {
	res := make(map[string]string)
	for k := range urls {
		res[k] = "http://short.ly/" + k
	}
	return res, nil
}

func ExampleSaveHandler_PostHandler() {
	h := handler.NewSaveHandler(mockSaver{})
	r := chi.NewRouter()
	r.Post("/", h.PostHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	reqBody := bytes.NewBufferString("https://example.com")
	resp, _ := http.Post(ts.URL+"/", "text/plain", reqBody)
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 201
}

func ExampleSaveHandler_PostHandlerJSON() {
	h := handler.NewSaveHandler(mockSaver{})
	r := chi.NewRouter()
	r.Post("/api/shorten", h.PostHandlerJSON)

	ts := httptest.NewServer(r)
	defer ts.Close()

	reqData := domain.ShortenRequest{URL: "https://example.com"}
	reqBytes, _ := json.Marshal(reqData)

	resp, _ := http.Post(ts.URL+"/api/shorten", "application/json", bytes.NewBuffer(reqBytes))
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 201
}

func ExampleSaveHandler_PostHandlerBatch() {
	h := handler.NewSaveHandler(mockSaver{})
	r := chi.NewRouter()
	r.Post("/api/shorten/batch", h.PostHandlerBatch)

	ts := httptest.NewServer(r)
	defer ts.Close()

	batchReq := []handler.BatchRequest{
		{CorrelationID: "1", OriginalURL: "https://site1.com"},
		{CorrelationID: "2", OriginalURL: "https://site2.com"},
	}
	reqBytes, _ := json.Marshal(batchReq)

	resp, _ := http.Post(ts.URL+"/api/shorten/batch", "application/json", bytes.NewBuffer(reqBytes))
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 201
}
