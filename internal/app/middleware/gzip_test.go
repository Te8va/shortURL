package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Te8va/shortURL/internal/app/middleware"
)

func TestGzipHandle(t *testing.T) {
	echoHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	})

	tests := []struct {
		name           string
		reqHeaders     map[string]string
		reqBody        []byte
		wantBody       string
		wantGzipEncode bool
	}{
		{
			name:           "no accept-encoding gzip",
			reqHeaders:     map[string]string{"Content-Type": "application/json"},
			reqBody:        []byte("plain body"),
			wantBody:       "plain body",
			wantGzipEncode: false,
		},
		{
			name: "with accept-encoding gzip",
			reqHeaders: map[string]string{
				"Content-Type":    "application/json",
				"Accept-Encoding": "gzip, deflate",
			},
			reqBody:        []byte("gzip response body"),
			wantBody:       "gzip response body",
			wantGzipEncode: true,
		},
		{
			name: "content-encoding gzip decompress",
			reqHeaders: map[string]string{
				"Content-Type":     "application/json",
				"Content-Encoding": "gzip",
			},
			reqBody: func() []byte {
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				gz.Write([]byte("compressed input"))
				gz.Close()
				return buf.Bytes()
			}(),
			wantBody:       "compressed input",
			wantGzipEncode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middleware.GzipHandle(echoHandler)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.reqBody))
			for k, v := range tt.reqHeaders {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			var body []byte
			var err error
			var gzr *gzip.Reader

			if tt.wantGzipEncode {
				require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
				gzr, err = gzip.NewReader(res.Body)
				require.NoError(t, err)
				body, err = io.ReadAll(gzr)
				require.NoError(t, err)
				gzr.Close()
			} else {
				body, err = io.ReadAll(res.Body)
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantBody, string(body))
		})
	}
}
