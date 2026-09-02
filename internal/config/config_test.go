package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	for _, key := range []string{"PORT", "READ_TIMEOUT", "WRITE_TIMEOUT", "SHUTDOWN_TIMEOUT", "OTEL_EXPORTER_OTLP_ENDPOINT", "CLICKHOUSE_DSN"} {
		t.Setenv(key, "")
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Address != ":8080" || got.ReadTimeout != 5*time.Second || got.WriteTimeout != 15*time.Second {
		t.Fatalf("unexpected defaults: %+v", got)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("PORT", "9091")
	t.Setenv("READ_TIMEOUT", "750ms")
	t.Setenv("WRITE_TIMEOUT", "2s")
	t.Setenv("SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://collector:4317")
	t.Setenv("CLICKHOUSE_DSN", "clickhouse://db:9000/traceforge")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Address != ":9091" || got.ReadTimeout != 750*time.Millisecond || got.OTLPEndpoint == "" || got.ClickHouseDSN == "" {
		t.Fatalf("overrides not loaded: %+v", got)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	t.Run("port", func(t *testing.T) {
		t.Setenv("PORT", "api")
		if _, err := Load(); err == nil {
			t.Fatal("expected invalid port error")
		}
	})
	t.Run("duration", func(t *testing.T) {
		t.Setenv("PORT", "8080")
		t.Setenv("READ_TIMEOUT", "soon")
		if _, err := Load(); err == nil {
			t.Fatal("expected invalid duration error")
		}
	})
}
