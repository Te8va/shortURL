package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/Te8va/shortURL/internal/app/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupTestHandler(t *testing.T) (*gomock.Controller, *mocks.MockURLSaver, *mocks.MockURLGetter, *mocks.MockPinger, *mocks.MockURLDelete, *URLHandler) {
	ctrl := gomock.NewController(t)

	mockSaver := mocks.NewMockURLSaver(ctrl)
	mockGetter := mocks.NewMockURLGetter(ctrl)
	mockPinger := mocks.NewMockPinger(ctrl)
	mockDeleter := mocks.NewMockURLDelete(ctrl)

	testCfg := &config.Config{
		BaseURL:       "http://localhost:8080",
		ServerAddress: "localhost:8080",
	}

	handler := NewURLHandler(testCfg, mockSaver, mockGetter, mockPinger, mockDeleter)
	return ctrl, mockSaver, mockGetter, mockPinger, mockDeleter, handler
}

func TestPostHandler(t *testing.T) {
	ctrl, mockSaver, _, _, _, handler := setupTestHandler(t)
	defer ctrl.Finish()

	testCases := []struct {
		name        string
		contentType string
		body        string
		mockReturn  string
		wantCode    int
	}{
		{
			name:        "valid URL",
			contentType: "text/plain",
			body:        "http://example.com",
			mockReturn:  "http://localhost:8080/shortID",
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.wantCode == http.StatusCreated {
				mockSaver.EXPECT().Save(gomock.Any(), gomock.Any(), testCase.body).Return(testCase.mockReturn, nil).Times(1)
			}

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testCase.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", testCase.contentType)

			w := httptest.NewRecorder()
			handler.PostHandler(w, req)

			require.Equal(t, testCase.wantCode, w.Code)
			if testCase.wantCode == http.StatusCreated {
				require.Contains(t, w.Body.String(), "http://localhost:8080/shortID")
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	ctrl, _, mockGetter, _, _, handler := setupTestHandler(t)
	defer ctrl.Finish()

	baseURL := "http://localhost:8080"
	testID := "testID"
	fullURL := fmt.Sprintf("%s/%s", baseURL, testID)
	testURL := "http://example.com"

	mockGetter.EXPECT().Get(gomock.Any(), fullURL).Return(testURL, true, false).AnyTimes()
	mockGetter.EXPECT().Get(gomock.Any(), fmt.Sprintf("%s/%s", baseURL, "deletedID")).Return("", true, true).AnyTimes()
	mockGetter.EXPECT().Get(gomock.Any(), fmt.Sprintf("%s/%s", baseURL, "invalidID")).Return("", false, false).AnyTimes()

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
			wantCode:  http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/"+testCase.requestID, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			handler.GetHandler(w, req)

			require.Equal(t, testCase.wantCode, w.Code)
			if testCase.wantCode == http.StatusTemporaryRedirect {
				require.Equal(t, testCase.wantURL, w.Header().Get("Location"))
			}
		})
	}
}

func TestPostHandlerJSON(t *testing.T) {
	ctrl, mockSaver, _, _, _, handler := setupTestHandler(t)
	defer ctrl.Finish()

	testCases := []struct {
		name        string
		contentType string
		body        domain.ShortenRequest
		mockReturn  string
		mockErr     error
		wantCode    int
	}{
		{
			name:        "valid JSON",
			contentType: "application/json",
			body:        domain.ShortenRequest{URL: "http://example.com"},
			mockReturn:  "http://localhost:8080/shortID",
			mockErr:     nil,
			wantCode:    http.StatusCreated,
		},
		{
			name:        "invalid content type",
			contentType: "text/plain",
			body:        domain.ShortenRequest{URL: "http://example.com"},
			wantCode:    http.StatusBadRequest,
		},
		{
			name:        "invalid URL",
			contentType: "application/json",
			body:        domain.ShortenRequest{URL: "invalid-url"},
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(testCase.body)

			if testCase.wantCode == http.StatusCreated {
				mockSaver.EXPECT().Save(gomock.Any(), gomock.Any(), testCase.body.URL).Return(testCase.mockReturn, testCase.mockErr).Times(1)
			}

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", testCase.contentType)

			w := httptest.NewRecorder()
			handler.PostHandlerJSON(w, req)

			require.Equal(t, testCase.wantCode, w.Code)
			if testCase.wantCode == http.StatusCreated {
				require.Contains(t, w.Body.String(), "http://localhost:8080/shortID")
			}
		})
	}
}

func TestPingHandler(t *testing.T) {
	ctrl, _, _, mockPinger, _, handler := setupTestHandler(t)
	defer ctrl.Finish()

	mockPinger.EXPECT().PingPg(gomock.Any()).Return(nil).Times(1)

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	handler.PingHandler(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
