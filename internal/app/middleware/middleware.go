package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

// Initialize sets up the global logger with the given log level
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging wraps an HTTP handler to log incoming requests and their responses
func WithLogging(h http.Handler) http.Handler {
	return GzipHandle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		duration := time.Since(start)

		Log.Info("Incoming request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
		)

		responseData := &responseData{
			status: http.StatusOK,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		Log.Info("HTTP request",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}))
}
