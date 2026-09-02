package ports

import (
	"context"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
)

type TraceFilter struct {
	Service     string
	Status      domain.SpanStatus
	MinDuration time.Duration
	From        time.Time
	To          time.Time
	Limit       int
}

type Repository interface {
	ListServices(context.Context) ([]domain.Service, error)
	ListTraces(context.Context, TraceFilter) ([]domain.Trace, error)
	SaveTrace(context.Context, domain.Trace) error
	ListSLOs(context.Context) ([]domain.SLO, error)
	ListDeployments(context.Context) ([]domain.Deployment, error)
	ListAlerts(context.Context) ([]domain.Alert, error)
	SaveAlert(context.Context, domain.Alert) error
	ListIncidents(context.Context) ([]domain.Incident, error)
	SaveIncident(context.Context, domain.Incident) error
	ListSuppressions(context.Context) ([]domain.SuppressionRule, error)
}
