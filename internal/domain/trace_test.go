package domain

import (
	"testing"
	"time"
)

func TestBuildTraceDerivesCriticalPathAndError(t *testing.T) {
	started := time.Date(2026, 9, 2, 8, 0, 0, 0, time.UTC)
	trace, err := BuildTrace("trace-1", []Span{{TraceID: "trace-1", SpanID: "root", Service: "gateway", Operation: "GET /", StartedAt: started, Duration: 200 * time.Millisecond, Status: SpanOK}, {TraceID: "trace-1", SpanID: "child", ParentID: "root", Service: "payments", Operation: "Authorize", StartedAt: started.Add(50 * time.Millisecond), Duration: 300 * time.Millisecond, Status: SpanError}})
	if err != nil {
		t.Fatal(err)
	}
	if !trace.HasError {
		t.Fatal("expected trace error")
	}
	if trace.RootService != "gateway" {
		t.Fatalf("unexpected root %q", trace.RootService)
	}
	if trace.Duration != 350*time.Millisecond {
		t.Fatalf("unexpected duration %s", trace.Duration)
	}
}
func TestBuildTraceRejectsMixedTraceIDs(t *testing.T) {
	_, err := BuildTrace("trace-1", []Span{{TraceID: "trace-2", SpanID: "span", Service: "api", Operation: "work", Status: SpanOK}})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
