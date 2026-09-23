package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", urlshortener)

	mux.HandleFunc("GET /{shortCode}", redirectUrl)

	mux.HandleFunc("GET /stats/{shortCode}", clickEventChecker)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go clickRecorder(ctx, eventClick)

	go requestLogCleaner(ctx)

	fmt.Println("Server listening on :8080")

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err.Error())
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan // blocks here until Ctrl+C

	slog.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err.Error())
	} else {
		slog.Info("server shut down cleanly")
	}

	cancel() // stops your goroutines

}
