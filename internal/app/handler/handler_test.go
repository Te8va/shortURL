package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Te8va/shortURL/internal/app/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCfg := &config.Config{
		BaseURL:       "http://localhost:8080",
		ServerAddress: "localhost:8080",
	}

	testStore := NewURLStore(testCfg)

	testCases := []struct {
		name        string
		contentType string
		body        string
		wantCode    int
	}{
		{
			name:        "valid URL",
			contentType: "text/plain",
			body:        "http://example.com",
			wantCode:    http.StatusCreated,
		},
		{
			name:        "invalid content type",
			contentType: "application/json",
			body:        "http://example.com",
			wantCode:    http.StatusBadRequest,
		},
		{
			name:        "empty URL",
			contentType: "text/plain",
			body:        "",
			wantCode:    http.StatusBadRequest,
		},
		{
			name:        "invalid URL format",
			contentType: "text/plain",
			body:        "invalid-url",
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testCase.body))
			require.NoError(t, err)

			req.Header.Set("Content-Type", testCase.contentType)

			w := httptest.NewRecorder()
			testStore.PostHandler(w, req)

			require.Equal(t, testCase.wantCode, w.Code)

			if testCase.wantCode == http.StatusBadRequest {
				require.NotEmpty(t, w.Body.String())
			} else {
				require.Contains(t, w.Body.String(), "http://localhost:8080/")
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCfg := &config.Config{
		BaseURL:       "http://localhost:8080",
		ServerAddress: "localhost:8080",
	}

	testStore := NewURLStore(testCfg)

	testID := "testID"
	testURL := "http://example.com"

	testStore.store[testID] = testURL

	testCases := []struct {
		name      string
		requestID string
		wantCode  int
		wantURL   string
	}{
		{
			name:      "valid ID",
			requestID: testID,
			wantCode:  http.StatusTemporaryRedirect,
			wantURL:   testURL,
		},
		{
			name:      "invalid ID",
			requestID: "invalidID",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "missing ID",
			requestID: "",
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/"+testCase.requestID, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			testStore.GetHandler(w, req)

			require.Equal(t, testCase.wantCode, w.Code)

			if testCase.wantCode == http.StatusTemporaryRedirect {
				require.Equal(t, testCase.wantURL, w.Header().Get("Location"))
			} else {
				require.NotEmpty(t, w.Body.String())
			}
		})
	}
}
