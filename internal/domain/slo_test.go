package domain

import "testing"

func TestSLOBudgetAndBurnRate(t *testing.T) {
	slo := SLO{ID: "slo", Service: "api", Name: "availability", Target: 99, GoodEvents: 985, TotalEvents: 1000}
	if got := slo.ErrorBudgetRemaining(); got != 0 {
		t.Fatalf("remaining budget = %v", got)
	}
	if got := slo.BurnRate(); got < 1.49 || got > 1.51 {
		t.Fatalf("burn rate = %v", got)
	}
}
func TestSLOWithoutTrafficIsHealthy(t *testing.T) {
	slo := SLO{Target: 99.9}
	if slo.Attainment() != 100 || slo.BurnRate() != 0 {
		t.Fatal("empty window should be healthy")
	}
}
