package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	config, err := uptime.LoadConfig()
	if err != nil {
		return err
	}
	handler, err := uptime.NewHandler()
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              config.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return fmt.Errorf("listen on UPTIME_ADDR %q: %w", config.Addr, err)
	}
	slog.Info("starting HTTP server", "address", listener.Addr().String())
	if err := serve(ctx, server, listener, 10*time.Second); err != nil {
		return err
	}
	slog.Info("HTTP server stopped")
	return nil
}
