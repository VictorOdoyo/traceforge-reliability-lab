package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type Repository struct {
	mu           sync.RWMutex
	Services     []domain.Service
	Traces       []domain.Trace
	SLOs         []domain.SLO
	Deployments  []domain.Deployment
	Alerts       []domain.Alert
	Incidents    []domain.Incident
	Suppressions []domain.SuppressionRule
}

func (r *Repository) ListServices(context.Context) ([]domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.Service(nil), r.Services...), nil
}
func (r *Repository) ListSLOs(context.Context) ([]domain.SLO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.SLO(nil), r.SLOs...), nil
}
func (r *Repository) ListDeployments(context.Context) ([]domain.Deployment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.Deployment(nil), r.Deployments...), nil
}
func (r *Repository) ListAlerts(context.Context) ([]domain.Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.Alert(nil), r.Alerts...), nil
}
func (r *Repository) ListIncidents(context.Context) ([]domain.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.Incident(nil), r.Incidents...), nil
}
func (r *Repository) ListSuppressions(context.Context) ([]domain.SuppressionRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.SuppressionRule(nil), r.Suppressions...), nil
}

func (r *Repository) ListTraces(_ context.Context, filter ports.TraceFilter) ([]domain.Trace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Trace, 0, len(r.Traces))
	for _, trace := range r.Traces {
		if filter.Service != "" && !traceContainsService(trace, filter.Service) {
			continue
		}
		if filter.Status == domain.SpanError && !trace.HasError {
			continue
		}
		if filter.Status == domain.SpanOK && trace.HasError {
			continue
		}
		if trace.Duration < filter.MinDuration {
			continue
		}
		if !filter.From.IsZero() && trace.StartedAt.Before(filter.From) {
			continue
		}
		if !filter.To.IsZero() && trace.StartedAt.After(filter.To) {
			continue
		}
		result = append(result, trace)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StartedAt.After(result[j].StartedAt) })
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}
	return result, nil
}

func (r *Repository) SaveTrace(_ context.Context, trace domain.Trace) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Traces = append(r.Traces, trace)
	return nil
}
func (r *Repository) SaveAlert(_ context.Context, alert domain.Alert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.Alerts {
		if r.Alerts[i].ID == alert.ID {
			r.Alerts[i] = alert
			return nil
		}
	}
	r.Alerts = append(r.Alerts, alert)
	return nil
}
func (r *Repository) SaveIncident(_ context.Context, incident domain.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.Incidents {
		if r.Incidents[i].ID == incident.ID {
			r.Incidents[i] = incident
			return nil
		}
	}
	r.Incidents = append(r.Incidents, incident)
	return nil
}

func traceContainsService(trace domain.Trace, service string) bool {
	for _, span := range trace.Spans {
		if span.Service == service {
			return true
		}
	}
	return false
}
