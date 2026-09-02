package domain

import (
	"testing"
	"time"
)

func TestSuppressionMatchesScopeAndWindow(t *testing.T) {
	now := time.Now()
	rule := SuppressionRule{Service: "inventory", LabelKey: "region", LabelValue: "east", StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}
	alert := Alert{Service: "inventory", Labels: map[string]string{"region": "east"}}
	if !rule.Matches(alert, now) {
		t.Fatal("expected suppression match")
	}
	alert.Labels["region"] = "west"
	if rule.Matches(alert, now) {
		t.Fatal("unexpected suppression match")
	}
}
func TestSamplingPrioritizesErrorsAndSlowTraces(t *testing.T) {
	policy := SamplingPolicy{DefaultRate: .05, ErrorRate: 1, SlowTraceRate: .5, SlowThresholdMillis: 500}
	if got := policy.RateFor(Trace{HasError: true}); got != 1 {
		t.Fatalf("error rate = %v", got)
	}
	if got := policy.RateFor(Trace{Duration: 600 * time.Millisecond}); got != .5 {
		t.Fatalf("slow rate = %v", got)
	}
	if got := policy.RateFor(Trace{}); got != .05 {
		t.Fatalf("default rate = %v", got)
	}
}
