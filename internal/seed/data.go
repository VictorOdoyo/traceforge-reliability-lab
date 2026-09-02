package seed

import (
	"fmt"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/platform/memory"
)

func Repository(now time.Time) *memory.Repository {
	services := []domain.Service{
		{Name: "edge-gateway", Team: "Platform", Tier: 1, Environment: "production", Language: "Go", Labels: map[string]string{"region": "eu-west"}},
		{Name: "checkout-api", Team: "Commerce", Tier: 1, Environment: "production", Language: "Kotlin", Labels: map[string]string{"region": "eu-west"}},
		{Name: "payment-router", Team: "Payments", Tier: 1, Environment: "production", Language: "Rust", Labels: map[string]string{"region": "eu-west"}},
		{Name: "inventory-api", Team: "Fulfillment", Tier: 2, Environment: "production", Language: "TypeScript", Labels: map[string]string{"region": "eu-west"}},
		{Name: "notification-worker", Team: "Experience", Tier: 3, Environment: "production", Language: "Python", Labels: map[string]string{"region": "eu-west"}},
	}
	traces := make([]domain.Trace, 0, 36)
	for index := 0; index < 36; index++ {
		trace, _ := syntheticTrace(index, now.Add(-time.Duration(index)*4*time.Minute))
		traces = append(traces, trace)
	}
	return &memory.Repository{
		Services: services, Traces: traces,
		SLOs: []domain.SLO{
			{ID: "slo-checkout", Service: "checkout-api", Name: "Successful checkout", Target: 99.9, Window: 30 * 24 * time.Hour, GoodEvents: 99840, TotalEvents: 100000},
			{ID: "slo-payment", Service: "payment-router", Name: "Payment authorization", Target: 99.95, Window: 30 * 24 * time.Hour, GoodEvents: 99920, TotalEvents: 100000},
			{ID: "slo-inventory", Service: "inventory-api", Name: "Inventory availability", Target: 99.5, Window: 30 * 24 * time.Hour, GoodEvents: 99740, TotalEvents: 100000},
		},
		Deployments: []domain.Deployment{
			{ID: "dep-1842", Service: "payment-router", Version: "v2.18.4", Environment: "production", CommitSHA: "8bc9f17", DeployedBy: "release-bot", DeployedAt: now.Add(-52 * time.Minute), ChangeSummary: "Retry policy and issuer routing update"},
			{ID: "dep-1841", Service: "checkout-api", Version: "v4.9.1", Environment: "production", CommitSHA: "154ae20", DeployedBy: "amina", DeployedAt: now.Add(-4 * time.Hour), ChangeSummary: "Cart validation performance improvements"},
		},
		Alerts: []domain.Alert{
			{ID: "alert-301", Service: "payment-router", Title: "Authorization error rate above 2%", Severity: domain.SeverityCritical, State: domain.AlertFiring, StartedAt: now.Add(-31 * time.Minute), UpdatedAt: now.Add(-4 * time.Minute), Labels: map[string]string{"region": "eu-west", "signal": "errors"}},
			{ID: "alert-302", Service: "checkout-api", Title: "Checkout p95 latency above 800ms", Severity: domain.SeverityHigh, State: domain.AlertAcknowledged, Owner: "Victor", StartedAt: now.Add(-42 * time.Minute), UpdatedAt: now.Add(-25 * time.Minute), Labels: map[string]string{"region": "eu-west", "signal": "latency"}},
			{ID: "alert-303", Service: "inventory-api", Title: "Inventory cache saturation", Severity: domain.SeverityMedium, State: domain.AlertFiring, StartedAt: now.Add(-12 * time.Minute), UpdatedAt: now.Add(-2 * time.Minute), Labels: map[string]string{"region": "eu-west", "signal": "capacity"}},
		},
		Incidents: []domain.Incident{{ID: "inc-1042", Title: "Intermittent card authorization failures", Severity: domain.SeverityCritical, Status: domain.IncidentInvestigating, Commander: "Samira", AffectedServices: []string{"checkout-api", "payment-router"}, StartedAt: now.Add(-28 * time.Minute), Timeline: []domain.TimelineEvent{
			{ID: "evt-1", Kind: "detected", Message: "Error budget burn alert crossed the fast-burn threshold", Actor: "alert-router", OccurredAt: now.Add(-28 * time.Minute)},
			{ID: "evt-2", Kind: "deployment", Message: "Payment router deployment identified inside correlation window", Actor: "TraceForge", OccurredAt: now.Add(-22 * time.Minute)},
		}}},
		Suppressions: []domain.SuppressionRule{{ID: "sup-12", Name: "Inventory maintenance", Service: "inventory-api", StartsAt: now.Add(2 * time.Hour), EndsAt: now.Add(3 * time.Hour), Reason: "Planned cache node replacement"}},
	}
}

func syntheticTrace(index int, startedAt time.Time) (domain.Trace, error) {
	id := fmt.Sprintf("7f3c%028x", index+1)
	gatewayDuration := time.Duration(280+index*17) * time.Millisecond
	paymentError := index%7 == 0
	spans := []domain.Span{
		{TraceID: id, SpanID: fmt.Sprintf("root-%02d", index), Service: "edge-gateway", Operation: "POST /checkout", StartedAt: startedAt, Duration: gatewayDuration, Status: domain.SpanOK, Attributes: map[string]string{"http.method": "POST", "region": "eu-west"}},
		{TraceID: id, SpanID: fmt.Sprintf("checkout-%02d", index), ParentID: fmt.Sprintf("root-%02d", index), Service: "checkout-api", Operation: "CreateOrder", StartedAt: startedAt.Add(12 * time.Millisecond), Duration: gatewayDuration - 30*time.Millisecond, Status: domain.SpanOK, Attributes: map[string]string{"rpc.system": "grpc"}},
		{TraceID: id, SpanID: fmt.Sprintf("payment-%02d", index), ParentID: fmt.Sprintf("checkout-%02d", index), Service: "payment-router", Operation: "Authorize", StartedAt: startedAt.Add(42 * time.Millisecond), Duration: time.Duration(110+index*12) * time.Millisecond, Status: status(paymentError), Attributes: map[string]string{"issuer": "demo-bank"}},
		{TraceID: id, SpanID: fmt.Sprintf("inventory-%02d", index), ParentID: fmt.Sprintf("checkout-%02d", index), Service: "inventory-api", Operation: "Reserve", StartedAt: startedAt.Add(46 * time.Millisecond), Duration: time.Duration(38+index) * time.Millisecond, Status: domain.SpanOK, Attributes: map[string]string{"cache.hit": "true"}},
	}
	return domain.BuildTrace(id, spans)
}

func status(failed bool) domain.SpanStatus {
	if failed {
		return domain.SpanError
	}
	return domain.SpanOK
}
