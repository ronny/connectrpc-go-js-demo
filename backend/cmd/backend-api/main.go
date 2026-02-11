package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/internal/api"
	"example.com/internal/service"
)

const GracefulShutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "backend-api returned an unexpected error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	host := cmp.Or(os.Getenv("HOST"), "localhost")
	port := cmp.Or(os.Getenv("PORT"), "8080")
	listenAddr := fmt.Sprintf("%s:%s", host, port)

	svc := service.New()

	apiServer, err := api.NewServer(listenAddr, svc)
	if err != nil {
		return fmt.Errorf("api.NewServer: %w", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := apiServer.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("backend-api server returned an unexpected error",
				"error", err)
		}

		slog.Info("backend-api stopped")
	}()

	slog.Info("backend-api started",
		slog.String(
			"address", apiServer.Address(),
		),
	)

	sig := <-signals
	slog.Info("received signal, shutting down backend-api gracefully...",
		"signal", sig,
		"gracefulShutdownTimeout", GracefulShutdownTimeout,
	)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	err = apiServer.Shutdown(shutdownCtx)
	if err != nil {
		slog.Warn("ignored error during apiServer.Shutdown",
			"error", err)
	}

	slog.Info("backend-api gracefully shut down")

	return nil
}
