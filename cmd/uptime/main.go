package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/StastkaJan/uptime-monitor/internal/uptime"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	handler, err := uptime.NewHandler()
	if err != nil {
		return err
	}
	addr := os.Getenv("UPTIME_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	slog.Info("starting HTTP server", "address", addr)
	return server.ListenAndServe()
}
