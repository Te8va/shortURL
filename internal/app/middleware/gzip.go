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
        if r.Header.Get("Content-Encoding") == "gzip" {
            gz, err := gzip.NewReader(r.Body)
            if err != nil {
                http.Error(w, "Failed to decompress request", http.StatusBadRequest)
                return
            }
            defer gz.Close()
            r.Body = io.NopCloser(gz) 
        }

        if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            w.Header().Set("Content-Encoding", "gzip")
            gz := gzip.NewWriter(w)
            defer gz.Close()
            gw := gzipWriter{ResponseWriter: w, Writer: gz}
            h.ServeHTTP(gw, r)
            return
        }

        h.ServeHTTP(w, r)
    })
}
