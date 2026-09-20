package uptime

import (
	"fmt"
	"net/http"
)

// NewHandler constructs the application's routes from embedded assets.
func NewHandler() (http.Handler, error) {
	dashboard, err := dashboardHTML()
	if err != nil {
		return nil, fmt.Errorf("render dashboard: %w", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(dashboard)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.Handle("GET /static/style.css", http.FileServer(http.FS(assets)))
	return mux, nil
}
