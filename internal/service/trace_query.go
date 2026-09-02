package service

import (
	"context"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type TraceQuery struct{ repository ports.Repository }

func NewTraceQuery(repository ports.Repository) TraceQuery { return TraceQuery{repository: repository} }

func (q TraceQuery) Search(ctx context.Context, service, status string, minimumMillis float64, limit int) ([]domain.Trace, error) {
	filter := ports.TraceFilter{Service: service, MinDuration: time.Duration(minimumMillis * float64(time.Millisecond)), Limit: limit}
	switch status {
	case "error":
		filter.Status = domain.SpanError
	case "ok":
		filter.Status = domain.SpanOK
	}
	return q.repository.ListTraces(ctx, filter)
}
