package handlers

import (
	"main/internal/handlers/sse"

	"github.com/go-chi/chi/v5"
)

func Register(r *chi.Mux) {
	r.Get("/sse", sse.Handler)
}
