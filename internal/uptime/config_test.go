package uptime

import (
	"os"
	"strings"
	"testing"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"UPTIME_ADDR", "UPTIME_DB", "UPTIME_DEMO"} {
		t.Setenv(key, "") // Register restoration even when the original value was unset.
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	clearConfigEnv(t)
	got, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	want := Config{Addr: "127.0.0.1:8080", DB: "data/uptime.db", Demo: false}
	if got != want {
		t.Fatalf("config = %+v, want %+v", got, want)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	for _, addr := range []string{"localhost:9090", "0.0.0.0:8080", "127.0.0.1:1", "[::1]:65535"} {
		t.Run(addr, func(t *testing.T) {
			clearConfigEnv(t)
			t.Setenv("UPTIME_ADDR", addr)
			t.Setenv("UPTIME_DB", "custom/monitor.db")
			t.Setenv("UPTIME_DEMO", "true")
			got, err := LoadConfig()
			if err != nil {
				t.Fatal(err)
			}
			want := Config{Addr: addr, DB: "custom/monitor.db", Demo: true}
			if got != want {
				t.Fatalf("config = %+v, want %+v", got, want)
			}
		})
	}
}

func TestLoadConfigInvalidEnvironment(t *testing.T) {
	for _, tc := range []struct {
		name, key, value, hint string
	}{
		{"empty address", "UPTIME_ADDR", "", "host:port"},
		{"missing port", "UPTIME_ADDR", "127.0.0.1", "host:port"},
		{"missing host", "UPTIME_ADDR", ":8080", "host:port"},
		{"host whitespace", "UPTIME_ADDR", " localhost:8080", "host:port"},
		{"unbracketed IPv6", "UPTIME_ADDR", "::1:8080", "host:port"},
		{"empty port", "UPTIME_ADDR", "localhost:", "1 to 65535"},
		{"named port", "UPTIME_ADDR", "localhost:http", "1 to 65535"},
		{"zero port", "UPTIME_ADDR", "localhost:0", "1 to 65535"},
		{"negative port", "UPTIME_ADDR", "localhost:-1", "1 to 65535"},
		{"large port", "UPTIME_ADDR", "localhost:65536", "1 to 65535"},
		{"empty database", "UPTIME_DB", "", "file path"},
		{"blank database", "UPTIME_DB", " \t ", "file path"},
		{"empty demo", "UPTIME_DEMO", "", "true or false"},
		{"invalid demo", "UPTIME_DEMO", "yes", "true or false"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clearConfigEnv(t)
			t.Setenv(tc.key, tc.value)
			_, err := LoadConfig()
			if err == nil {
				t.Fatal("expected configuration error")
			}
			if !strings.Contains(err.Error(), tc.key) || !strings.Contains(err.Error(), tc.hint) {
				t.Errorf("error = %q, want variable %s and guidance %q", err, tc.key, tc.hint)
			}
		})
	}
}
