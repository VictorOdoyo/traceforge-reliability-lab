package service

import (
	"context"
	"sort"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type DependencyEdge struct {
	Source           string  `json:"source"`
	Target           string  `json:"target"`
	Calls            int     `json:"calls"`
	Errors           int     `json:"errors"`
	AverageLatencyMs float64 `json:"averageLatencyMs"`
}
type edgeAggregate struct {
	calls, errors int
	latency       float64
}

type TopologyService struct{ repository ports.Repository }

func NewTopologyService(repository ports.Repository) TopologyService {
	return TopologyService{repository: repository}
}

func (s TopologyService) Edges(ctx context.Context) ([]DependencyEdge, error) {
	traces, err := s.repository.ListTraces(ctx, ports.TraceFilter{})
	if err != nil {
		return nil, err
	}
	aggregates := map[string]*edgeAggregate{}
	for _, trace := range traces {
		spans := map[string]string{}
		for _, span := range trace.Spans {
			spans[span.SpanID] = span.Service
		}
		for _, span := range trace.Spans {
			source := spans[span.ParentID]
			if source == "" || source == span.Service {
				continue
			}
			key := source + "\x00" + span.Service
			aggregate := aggregates[key]
			if aggregate == nil {
				aggregate = &edgeAggregate{}
				aggregates[key] = aggregate
			}
			aggregate.calls++
			aggregate.latency += span.DurationMillis()
			if span.Status == "error" {
				aggregate.errors++
			}
		}
	}
	result := make([]DependencyEdge, 0, len(aggregates))
	for key, value := range aggregates {
		split := 0
		for i := range key {
			if key[i] == 0 {
				split = i
				break
			}
		}
		result = append(result, DependencyEdge{Source: key[:split], Target: key[split+1:], Calls: value.calls, Errors: value.errors, AverageLatencyMs: value.latency / float64(value.calls)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Calls > result[j].Calls })
	return result, nil
}
