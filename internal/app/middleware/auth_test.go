package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/Te8va/shortURL/internal/app/middleware"
	"github.com/stretchr/testify/require"
)

const secretKey = "test_secret"

func TestAuthMiddleware_TableDriven(t *testing.T) {
	token, validUserID, err := middleware.GenerateToken(secretKey)
	require.NoError(t, err)

	tests := []struct {
		name             string
		cookie           *http.Cookie
		expectNewCookie  bool
		expectStatusCode int
		expectUserID     int
	}{
		{
			name:             "no auth cookie",
			cookie:           nil,
			expectNewCookie:  true,
			expectStatusCode: http.StatusOK,
			expectUserID:     0,
		},
		{
			name: "valid auth cookie",
			cookie: &http.Cookie{
				Name:  "auth",
				Value: token,
			},
			expectNewCookie:  false,
			expectStatusCode: http.StatusOK,
			expectUserID:     validUserID,
		},
		{
			name: "invalid auth cookie format",
			cookie: &http.Cookie{
				Name:  "auth",
				Value: "broken.token.value",
			},
			expectNewCookie:  true,
			expectStatusCode: http.StatusOK,
			expectUserID:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUserID int

			handler := middleware.AuthMiddleware(secretKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				val := r.Context().Value(domain.UserIDKey)
				require.NotNil(t, val)
				capturedUserID = val.(int)
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tt.expectNewCookie {
				require.NotNil(t, findCookie(res.Cookies(), "auth"))
			} else {
				require.Nil(t, findCookie(res.Cookies(), "auth"))
			}

			require.Equal(t, tt.expectStatusCode, res.StatusCode)

			if tt.expectUserID != 0 {
				require.Equal(t, tt.expectUserID, capturedUserID)
			} else {
				require.NotZero(t, capturedUserID)
			}
		})
	}
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}
