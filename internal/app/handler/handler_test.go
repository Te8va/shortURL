package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Te8va/shortURL/internal/app/config"
	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/Te8va/shortURL/internal/app/handler/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupTestHandler(t *testing.T) (*gomock.Controller, *mocks.MockURLSaver, *mocks.MockURLGetter, *mocks.MockPinger, *SaveHandler, *GetterHandler, *PingHandler) {
	ctrl := gomock.NewController(t)

	mockSaver := mocks.NewMockURLSaver(ctrl)
	mockGetter := mocks.NewMockURLGetter(ctrl)
	mockPinger := mocks.NewMockPinger(ctrl)

	testCfg := &config.Config{
		BaseURL:       "http://localhost:8080",
		ServerAddress: "localhost:8080",
	}

	saveHandler := NewSaveHandler(mockSaver)
	getterHandler := NewGetterHandler(mockGetter, testCfg)
	pingHandler := NewPingHandler(mockPinger)

	return ctrl, mockSaver, mockGetter, mockPinger, saveHandler, getterHandler, pingHandler
}

func TestPostHandler(t *testing.T) {
	ctrl, mockSaver, _, _, saveHandler, _, _ := setupTestHandler(t)
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
			saveHandler.PostHandler(w, req)

			require.Equal(t, testCase.wantCode, w.Code)
			if testCase.wantCode == http.StatusCreated {
				require.Contains(t, w.Body.String(), "http://localhost:8080/shortID")
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	ctrl, _, mockGetter, _, _, getterHandler, _ := setupTestHandler(t)
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
			getterHandler.GetHandler(w, req)

			require.Equal(t, testCase.wantCode, w.Code)
			if testCase.wantCode == http.StatusTemporaryRedirect {
				require.Equal(t, testCase.wantURL, w.Header().Get("Location"))
			}
		})
	}
}

func TestPostHandlerJSON(t *testing.T) {
	ctrl, mockSaver, _, _, saveHandler, _, _ := setupTestHandler(t)
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
			saveHandler.PostHandlerJSON(w, req)

			require.Equal(t, testCase.wantCode, w.Code)
			if testCase.wantCode == http.StatusCreated {
				require.Contains(t, w.Body.String(), "http://localhost:8080/shortID")
			}
		})
	}
}

func TestPingHandler(t *testing.T) {
	ctrl, _, _, mockPinger, _, _, pingHandler := setupTestHandler(t)
	defer ctrl.Finish()

	mockPinger.EXPECT().PingPg(gomock.Any()).Return(nil).Times(1)

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	pingHandler.PingHandler(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUserURLsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeleter := mocks.NewMockURLDelete(ctrl)
	testCfg := &config.Config{BaseURL: "http://localhost:8080"}
	deleteHandler := NewDeleteHandler(mockDeleter, testCfg)

	testCases := []struct {
		name       string
		body       interface{}
		userID     interface{}
		wantCode   int
		expectCall bool
	}{
		{
			name:       "valid request",
			body:       []string{"abc123", "def456"},
			userID:     42,
			wantCode:   http.StatusAccepted,
			expectCall: true,
		},
		{
			name:       "unauthorized",
			body:       []string{"abc123"},
			userID:     nil,
			wantCode:   http.StatusUnauthorized,
			expectCall: false,
		},
		{
			name:       "invalid JSON",
			body:       "{not:json}",
			userID:     42,
			wantCode:   http.StatusBadRequest,
			expectCall: false,
		},
		{
			name:       "empty list",
			body:       []string{},
			userID:     42,
			wantCode:   http.StatusBadRequest,
			expectCall: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			switch v := tc.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(bodyBytes))
			if tc.userID != nil {
				req = req.WithContext(context.WithValue(req.Context(), domain.UserIDKey, tc.userID))
			}

			if tc.expectCall {
				var ids []string
				_ = json.Unmarshal(bodyBytes, &ids)
				var full []string
				for _, id := range ids {
					full = append(full, fmt.Sprintf("%s/%s", testCfg.BaseURL, id))
				}

				done := make(chan struct{})
				mockDeleter.EXPECT().
					DeleteUserURLs(gomock.Any(), full, tc.userID.(int)).
					DoAndReturn(func(ctx context.Context, ids []string, userID int) error {
						close(done)
						return nil
					}).Times(1)

				w := httptest.NewRecorder()
				deleteHandler.DeleteUserURLsHandler(w, req)

				<-done
				require.Equal(t, tc.wantCode, w.Code)
			} else {
				w := httptest.NewRecorder()
				deleteHandler.DeleteUserURLsHandler(w, req)
				require.Equal(t, tc.wantCode, w.Code)
			}
		})
	}

}
func TestGetUserURLsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mocks.NewMockURLGetter(ctrl)
	testCfg := &config.Config{BaseURL: "http://localhost:8080"}
	handler := NewGetterHandler(mockGetter, testCfg)

	testCases := []struct {
		name       string
		userID     interface{}
		mockResult []map[string]string
		mockErr    error
		wantCode   int
		wantBody   []map[string]string
	}{
		{
			name:   "authorized with URLs",
			userID: 123,
			mockResult: []map[string]string{
				{
					"short_url":    "http://localhost:8080/abc",
					"original_url": "http://example.com",
				},
			},
			mockErr:  nil,
			wantCode: http.StatusOK,
			wantBody: []map[string]string{
				{
					"short_url":    "http://localhost:8080/abc",
					"original_url": "http://example.com",
				},
			},
		},
		{
			name:       "authorized with no URLs",
			userID:     123,
			mockResult: nil,
			mockErr:    nil,
			wantCode:   http.StatusNoContent,
		},
		{
			name:     "unauthorized",
			userID:   nil,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "internal error",
			userID:   123,
			mockErr:  fmt.Errorf("db error"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			if tc.userID != nil {
				req = req.WithContext(context.WithValue(req.Context(), domain.UserIDKey, tc.userID))
			}

			if tc.userID != nil {
				mockGetter.EXPECT().
					GetUserURLs(gomock.Any(), tc.userID.(int)).
					Return(tc.mockResult, tc.mockErr).
					Times(1)
			}

			w := httptest.NewRecorder()
			handler.GetUserURLsHandler(w, req)

			require.Equal(t, tc.wantCode, w.Code)

			if tc.wantCode == http.StatusOK {
				var got []map[string]string
				err := json.NewDecoder(w.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tc.wantBody, got)
			}
		})
	}
}
