package config

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)


type Config struct {
	TrivyCacheDir string
	TrivyTimeout  time.Duration
	ScanInterval  time.Duration
	ScanWorkers   int
	ScanStagger   time.Duration
	ListenAddr    string
	LogLevel      slog.Level
	LogFormat     string
}

func Load(getenv func(string) string) (*Config, error) {
	trivyCacheDir := envOrDefault(getenv, "VORDR_TRIVY_CACHE_DIR", "/var/cache/trivy")

	trivyTimeout, err := parseDuration(getenv, "VORDR_TRIVY_TIMEOUT", "5m")
	if err != nil {
		return nil, err
	}

	scanInterval, err := parseDuration(getenv, "VORDR_SCAN_INTERVAL", "6h")
	if err != nil {
		return nil, err
	}

	scanWorkers, err := parseInt(getenv, "VORDR_SCAN_WORKERS", "3")
	if err != nil {
		return nil, err
	}

	scanStagger, err := parseDuration(getenv, "VORDR_SCAN_STAGGER", "30s")
	if err != nil {
		return nil, err
	}

	listenAddr := envOrDefault(getenv, "VORDR_LISTEN_ADDR", ":9550")
	if _, _, err := net.SplitHostPort(listenAddr); err != nil {
		return nil, fmt.Errorf("invalid VORDR_LISTEN_ADDR %q: %w", listenAddr, err)
	}

	logLevel, err := parseLogLevel(envOrDefault(getenv, "VORDR_LOG_LEVEL", "info"))
	if err != nil {
		return nil, err
	}

	logFormat := envOrDefault(getenv, "VORDR_LOG_FORMAT", "text")
	if logFormat != "text" && logFormat != "json" {
		return nil, fmt.Errorf("VORDR_LOG_FORMAT must be 'text' or 'json', got %q", logFormat)
	}

	return &Config{
		TrivyCacheDir: trivyCacheDir,
		TrivyTimeout:  trivyTimeout,
		ScanInterval:  scanInterval,
		ScanWorkers:  scanWorkers,
		ScanStagger:  scanStagger,
		ListenAddr:   listenAddr,
		LogLevel:     logLevel,
		LogFormat:    logFormat,
	}, nil
}

func envOrDefault(getenv func(string) string, key, defaultVal string) string {
	if v := getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseDuration(getenv func(string) string, envKey, defaultVal string) (time.Duration, error) {
	raw := envOrDefault(getenv, envKey, defaultVal)
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", envKey, raw, err)
	}
	return d, nil
}

func parseInt(getenv func(string) string, envKey, defaultVal string) (int, error) {
	raw := envOrDefault(getenv, envKey, defaultVal)
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", envKey, raw, err)
	}
	return n, nil
}

func parseLogLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid VORDR_LOG_LEVEL %q: must be debug, info, warn, or error", s)
	}
}
