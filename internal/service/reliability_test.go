package service

import (
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
)

func TestSummarizeServiceComputesErrorAndLatencySignals(t *testing.T) {
	traces := []domain.Trace{{Spans: []domain.Span{{Service: "api", Duration: 100 * time.Millisecond, Status: domain.SpanOK}}}, {Spans: []domain.Span{{Service: "api", Duration: 900 * time.Millisecond, Status: domain.SpanError}}}}
	summary := summarizeService("api", traces)
	if summary.TraceCount != 2 || summary.ErrorCount != 1 {
		t.Fatalf("unexpected counts: %#v", summary)
	}
	if summary.ErrorRate != 50 || summary.AverageLatencyMs != 500 || summary.P95LatencyMs != 900 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
	if summary.HealthScore >= 100 {
		t.Fatal("degraded service should lose health score")
	}
}
