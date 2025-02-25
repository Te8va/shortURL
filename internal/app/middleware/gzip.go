package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func gzipHandle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" && r.Header.Get("Content-Type") != "text/html" && r.Header.Get("Content-Type") != "application/x-gzip" {
			h.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request", http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			r.Body= gz
			r.Header.Set("Content-Type", "text/plain")
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
