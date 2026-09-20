package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServeDrainsActiveRequest(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	stopping := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		_, _ = io.WriteString(w, "finished")
	})}
	server.RegisterOnShutdown(func() { close(stopping) })
	listener := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- serve(ctx, server, listener, time.Second) }()
	response := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		res, err := client.Get("http://" + listener.Addr().String())
		if err == nil {
			defer res.Body.Close()
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil && string(body) != "finished" {
				err = errors.New("active request did not finish")
			}
		}
		response <- err
	}()
	waitSignal(t, entered)
	cancel()
	waitSignal(t, stopping)
	select {
	case err := <-done:
		t.Fatalf("server stopped before active request finished: %v", err)
	default:
	}
	release <- struct{}{}
	if err := waitResult(t, response); err != nil {
		t.Fatal(err)
	}
	if err := waitResult(t, done); err != nil {
		t.Fatal(err)
	}
}

func TestServeClosesRequestsAfterShutdownDeadline(t *testing.T) {
	entered := make(chan struct{})
	canceled := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-r.Context().Done()
		close(canceled)
	})}
	listener := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- serve(ctx, server, listener, 50*time.Millisecond) }()
	response := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		res, err := client.Get("http://" + listener.Addr().String())
		if res != nil {
			res.Body.Close()
		}
		response <- err
	}()
	waitSignal(t, entered)
	cancel()
	if err := waitResult(t, done); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("shutdown error = %v, want deadline exceeded", err)
	}
	waitSignal(t, canceled)
	if err := waitResult(t, response); err == nil {
		t.Fatal("request unexpectedly completed successfully")
	}
}

func TestServeReportsListenerFailure(t *testing.T) {
	listener := listen(t)
	listener.Close()
	err := serve(context.Background(), &http.Server{}, listener, time.Second)
	if err == nil || !strings.Contains(err.Error(), "serve HTTP") {
		t.Fatalf("error = %v, want actionable serve failure", err)
	}
}

func listen(t *testing.T) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	return listener
}

func waitSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for lifecycle event")
	}
}

func waitResult(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for lifecycle result")
		return nil
	}
}
