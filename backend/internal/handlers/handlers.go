package handlers

import (
	"main/internal/handlers/health"
	"main/internal/handlers/sse"

	"github.com/go-chi/chi/v5"
)

func Register(r *chi.Mux) {
	r.Get("/health", health.Handler)

	r.Get("/sse", sse.Handler)
}
