package sse

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"main/internal/logger"
	"main/internal/middleware"
	"main/internal/quote"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// check if the request supports sse
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusBadRequest)
		return
	}

	// set up connection for sse
	w.Header().Set("Content-Type", "text/event-stream; charset=UTF-8")
	w.Header().Set("Cache-Control", "No-Cache")
	w.Header().Set("Connection", "Keep-Alive")

	// send headers to client
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	// setup logger
	eventLogger := logger.Stdout.With(slog.String("request_id", middleware.GetReqID(r.Context())))

	var eventCount uint64

	// convenient function to call at connection and in a loop
	sendData := func() {
		fmt.Fprintf(w, "data: %s\n\n", quote.GetRandom())

		flusher.Flush()

		eventCount++
	}

	// log sse request
	eventLogger.Info("sse client connected")

	// send data on connection
	sendData()

	// setup scheduler to send data every 1 second
	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	done := make(chan struct{}, 1)

	// wait for client to disconnect then close the done channel
	go func() {
		<-r.Context().Done()

		eventLogger.Info("sse client disconnected", slog.Uint64("event_count", eventCount))

		close(done)
	}()

	// wait for tick or done channel, send data every tick
	for {
		select {
		case <-ticker.C:
			sendData()
		case <-done:
			return
		}
	}
}
