package seed

import (
	"context"
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

func TestRepositoryProvidesConnectedDemoDataset(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	repository := Repository(now)

	services, _ := repository.ListServices(context.Background())
	traces, _ := repository.ListTraces(context.Background(), ports.TraceFilter{})
	slos, _ := repository.ListSLOs(context.Background())
	alerts, _ := repository.ListAlerts(context.Background())
	incidents, _ := repository.ListIncidents(context.Background())
	if len(services) != 5 || len(traces) != 36 || len(slos) != 3 || len(alerts) != 3 || len(incidents) != 1 {
		t.Fatalf("incomplete dataset: services=%d traces=%d slos=%d alerts=%d incidents=%d", len(services), len(traces), len(slos), len(alerts), len(incidents))
	}
	if traces[0].StartedAt.After(now) || incidents[0].AffectedServices[0] != "checkout-api" {
		t.Fatal("seed relationships or timestamps are invalid")
	}
}
