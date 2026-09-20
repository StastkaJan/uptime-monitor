package uptime

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr string
	DB   string
	Demo bool
}

func LoadConfig() (Config, error) {
	config := Config{Addr: "127.0.0.1:8080", DB: "data/uptime.db"}
	if value, ok := os.LookupEnv("UPTIME_ADDR"); ok {
		config.Addr = value
	}
	host, port, err := net.SplitHostPort(config.Addr)
	if err != nil || host == "" || strings.TrimSpace(host) != host {
		return Config{}, fmt.Errorf("UPTIME_ADDR must be host:port, for example 127.0.0.1:8080")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("UPTIME_ADDR port must be a number from 1 to 65535")
	}
	if value, ok := os.LookupEnv("UPTIME_DB"); ok {
		config.DB = value
	}
	if strings.TrimSpace(config.DB) == "" || strings.ContainsRune(config.DB, 0) {
		return Config{}, fmt.Errorf("UPTIME_DB must be a nonempty file path")
	}
	if value, ok := os.LookupEnv("UPTIME_DEMO"); ok {
		config.Demo, err = strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("UPTIME_DEMO must be a boolean, for example true or false")
		}
	}
	return config, nil
}
