package telemetry

import (
	"context"
	"testing"
)

func TestSetupWithoutExporterProvidesCleanShutdown(t *testing.T) {
	shutdown, err := Setup(context.Background(), "")
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected shutdown function")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error = %v", err)
	}
}
