package main

import (
	"cmp"
	"log/slog"
	"net/http"
	"os"
	"time"

	"main/internal/handlers"
	"main/internal/logger"
	"main/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	middleware.Register(r)

	handlers.Register(r)

	port := cmp.Or(os.Getenv("PORT"), "3000")

	var s = &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		ReadHeaderTimeout: 1 * time.Second,
	}

	logger.Stdout.Info("starting server", slog.String("port", port))

	if err := s.ListenAndServe(); err != nil {
		logger.Stdout.Info("server exited", logger.ErrAttr(err))
	}
}
