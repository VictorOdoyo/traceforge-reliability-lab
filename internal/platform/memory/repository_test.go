package memory

import (
	"context"
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

func TestListTracesFiltersAndLimits(t *testing.T) {
	now := time.Now()
	repository := &Repository{Traces: []domain.Trace{{ID: "old", RootService: "api", StartedAt: now.Add(-time.Hour), Duration: 100 * time.Millisecond}, {ID: "recent", RootService: "api", StartedAt: now, Duration: 800 * time.Millisecond, HasError: true, Spans: []domain.Span{{Service: "payments"}}}}}
	traces, err := repository.ListTraces(context.Background(), ports.TraceFilter{Service: "payments", Status: domain.SpanError, MinDuration: 500 * time.Millisecond, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(traces) != 1 || traces[0].ID != "recent" {
		t.Fatalf("unexpected traces: %#v", traces)
	}
}
func TestSaveAlertUpdatesExistingIdentity(t *testing.T) {
	repository := &Repository{Alerts: []domain.Alert{{ID: "alert", State: domain.AlertFiring}}}
	if err := repository.SaveAlert(context.Background(), domain.Alert{ID: "alert", State: domain.AlertAcknowledged}); err != nil {
		t.Fatal(err)
	}
	if len(repository.Alerts) != 1 || repository.Alerts[0].State != domain.AlertAcknowledged {
		t.Fatal("alert was not updated")
	}
}
