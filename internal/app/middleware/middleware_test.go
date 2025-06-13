package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/stretchr/testify/require"
)

func TestWithLogging(t *testing.T) {
	err := middleware.Initialize("error")
	require.NoError(t, err)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) 
		_, _ = w.Write([]byte("hello world"))
	})

	handler := middleware.WithLogging(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test/uri", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	
	body := w.Body.String()

	require.Equal(t, http.StatusTeapot, resp.StatusCode)
	require.Equal(t, "hello world", body)
}
