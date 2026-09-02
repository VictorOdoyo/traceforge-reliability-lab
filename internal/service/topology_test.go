package service

import (
	"context"
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/platform/memory"
)

func TestTopologyAggregatesCrossServiceSpans(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	trace, err := domain.BuildTrace("trace-1", []domain.Span{
		{TraceID: "trace-1", SpanID: "root", Service: "gateway", Operation: "POST /checkout", StartedAt: now, Duration: 200 * time.Millisecond, Status: domain.SpanOK},
		{TraceID: "trace-1", SpanID: "child", ParentID: "root", Service: "checkout", Operation: "CreateOrder", StartedAt: now, Duration: 80 * time.Millisecond, Status: domain.SpanError},
		{TraceID: "trace-1", SpanID: "nested", ParentID: "child", Service: "checkout", Operation: "validate", StartedAt: now, Duration: 10 * time.Millisecond, Status: domain.SpanOK},
	})
	if err != nil {
		t.Fatal(err)
	}
	repository := &memory.Repository{Traces: []domain.Trace{trace}}

	edges, err := NewTopologyService(repository).Edges(context.Background())
	if err != nil {
		t.Fatalf("Edges() error = %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("len(edges) = %d, want 1", len(edges))
	}
	if edges[0].Source != "gateway" || edges[0].Target != "checkout" || edges[0].Errors != 1 || edges[0].AverageLatencyMs != 80 {
		t.Fatalf("unexpected edge: %+v", edges[0])
	}
}
