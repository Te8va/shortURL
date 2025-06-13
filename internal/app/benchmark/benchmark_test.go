package benchmark

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/router"
	"github.com/Te8va/shortURL/internal/app/service/mocks"
)

var (
	r    chi.Router
	once sync.Once
)

func initBenchmarkRouter(b *testing.B) {
	once.Do(func() {
		ctrl := gomock.NewController(b)

		mockSaver := mocks.NewMockURLSaverServ(ctrl)
		mockGetter := mocks.NewMockURLGetterServ(ctrl)
		mockDeleter := mocks.NewMockURLDeleteServ(ctrl)

		mockSaver.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return("shortURL", nil).AnyTimes()
		mockSaver.EXPECT().SaveBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockGetter.EXPECT().Get(gomock.Any(), gomock.Any()).Return("https://example.com",true, true).AnyTimes()
		mockGetter.EXPECT().GetUserURLs(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockDeleter.EXPECT().DeleteUserURLs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		cfg := config.NewConfig()
		r = router.NewRouter(cfg, mockSaver, mockGetter, nil, mockDeleter)
	})
}

func BenchmarkPostShorten(b *testing.B) {
	initBenchmarkRouter(b)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
		req.Header.Set("Content-Type", "text/plain")
		req.AddCookie(&http.Cookie{Name: "token", Value: "mock-token"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGetOriginalURL(b *testing.B) {
	initBenchmarkRouter(b)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
		req.AddCookie(&http.Cookie{Name: "token", Value: "mock-token"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkBatchShorten(b *testing.B) {
	initBenchmarkRouter(b)

	body := `
[
	{"correlation_id": "1", "original_url": "https://example.com/1"},
	{"correlation_id": "2", "original_url": "https://example.com/2"}
]`

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "token", Value: "mock-token"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGetUserURLs(b *testing.B) {
	initBenchmarkRouter(b)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req.AddCookie(&http.Cookie{Name: "token", Value: "mock-token"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkDeleteUserURLs(b *testing.B) {
	initBenchmarkRouter(b)

	body := `["abc123", "xyz789"]`

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "token", Value: "mock-token"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
