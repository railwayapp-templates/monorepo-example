package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"main/internal/logger"
)

// logger middleware for access logs
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// gathers metrics from the upstream handlers
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			requestStartTime := time.Now()

			next.ServeHTTP(ww, r)

			requestEndTime := time.Since(requestStartTime)

			//prints log and metrics
			logger.Stdout.Info("handled request",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.RequestURI()),
				slog.String("hostname", r.Host),
				slog.String("user_agent", r.Header.Get("User-Agent")),
				slog.String("ip", r.RemoteAddr),
				slog.Int("code", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("request_id", GetReqID(r.Context())),
				logger.DurationAttr(requestEndTime, "request_time"),
			)
		})
	}
}
