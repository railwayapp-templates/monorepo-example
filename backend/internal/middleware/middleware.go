package middleware

import (
	"main/internal/config"
	"main/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Register(r *chi.Mux) {
	r.Use(middleware.Recoverer)

	r.Use(TrustProxy(&trustProxyConfig{
		ErrorLogger: logger.Stderr,
	}))

	r.Use(RequestID)

	r.Use(Logger())

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: config.Cors.AllowedOrigins,
	}))
}
