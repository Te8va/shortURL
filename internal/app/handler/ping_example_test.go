package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/Te8va/shortURL/internal/app/handler"
)

type mockPinger struct{}

func (m mockPinger) PingPg(ctx context.Context) error {
	return nil
}

func ExamplePingHandler_PingHandler() {
	ph := handler.NewPingHandler(mockPinger{})

	r := chi.NewRouter()
	r.Get("/ping", ph.PingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := http.Get(ts.URL + "/ping")
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Output:
	// 200
}
