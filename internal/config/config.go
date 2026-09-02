package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	OTLPEndpoint    string
	ClickHouseDSN   string
}

func Load() (Config, error) {
	port := env("PORT", "8080")
	if _, err := strconv.Atoi(port); err != nil {
		return Config{}, fmt.Errorf("PORT must be numeric: %w", err)
	}
	read, err := duration("READ_TIMEOUT", "5s")
	if err != nil {
		return Config{}, err
	}
	write, err := duration("WRITE_TIMEOUT", "15s")
	if err != nil {
		return Config{}, err
	}
	shutdown, err := duration("SHUTDOWN_TIMEOUT", "10s")
	if err != nil {
		return Config{}, err
	}
	return Config{Address: ":" + port, ReadTimeout: read, WriteTimeout: write, ShutdownTimeout: shutdown, OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), ClickHouseDSN: os.Getenv("CLICKHOUSE_DSN")}, nil
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func duration(key, fallback string) (time.Duration, error) {
	value := env(key, fallback)
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration: %w", key, err)
	}
	return parsed, nil
}
