package uptime

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestDashboardAndStylesheet(t *testing.T) {
	handler, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	page := httptest.NewRecorder()
	handler.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d, want 200", page.Code)
	}
	if got := page.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Errorf("dashboard content type = %q, want text/html", got)
	}
	body := page.Body.String()
	for _, fragment := range []string{"<!doctype html>", "<html", "<head>", "<title>uptime monitor</title>", "<body>", "<main", "<h1", "</html>"} {
		if !strings.Contains(strings.ToLower(body), fragment) {
			t.Errorf("dashboard is missing %q", fragment)
		}
	}
	link := regexp.MustCompile(`<link\b[^>]*\bhref="([^"]+\.css)"[^>]*>`).FindStringSubmatch(body)
	if len(link) != 2 {
		t.Fatal("dashboard has no CSS link")
	}
	if link[1] != "/static/style.css" || !strings.Contains(link[0], `rel="stylesheet"`) {
		t.Fatalf("stylesheet link = %q, want local /static/style.css stylesheet", link[0])
	}
	css := httptest.NewRecorder()
	handler.ServeHTTP(css, httptest.NewRequest(http.MethodGet, link[1], nil))
	if css.Code != http.StatusOK {
		t.Fatalf("stylesheet status = %d, want 200", css.Code)
	}
	if got := css.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/css") {
		t.Errorf("stylesheet content type = %q, want text/css", got)
	}
	if body := css.Body.String(); !strings.Contains(body, "{") || !strings.Contains(body, "}") {
		t.Error("stylesheet contains no CSS rule")
	}
}

func TestHealth(t *testing.T) {
	handler, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusOK {
		t.Errorf("health status = %d, want 200", response.Code)
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
		t.Errorf("health content type = %q, want text/plain", got)
	}
	if got := response.Body.String(); got != "ok\n" {
		t.Errorf("health body = %q, want ok followed by newline", got)
	}
}

func TestRouteBoundaries(t *testing.T) {
	handler, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/missing", http.StatusNotFound},
		{http.MethodGet, "/healthz/missing", http.StatusNotFound},
		{http.MethodGet, "/static/missing.css", http.StatusNotFound},
		{http.MethodPost, "/", http.StatusMethodNotAllowed},
		{http.MethodPost, "/healthz", http.StatusMethodNotAllowed},
		{http.MethodPost, "/static/style.css", http.StatusMethodNotAllowed},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, nil))
			if response.Code != tc.status {
				t.Errorf("status = %d, want %d", response.Code, tc.status)
			}
		})
	}
}
