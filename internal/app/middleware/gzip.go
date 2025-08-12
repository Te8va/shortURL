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

// Write writes the data to the gzip stream, wrapping the original http.ResponseWriter
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandle wraps an HTTP handler to transparently handle gzip compression and decompression for requests and responses when supported.
func GzipHandle(h http.Handler) http.Handler {
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

			r.Body = gz
			r.Header.Set("Content-Type", "text/plain")
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			if _, writeErr := io.WriteString(w, err.Error()); writeErr != nil {
				http.Error(w, "Failed to write error message", http.StatusInternalServerError)
			}
			return
		}

		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
