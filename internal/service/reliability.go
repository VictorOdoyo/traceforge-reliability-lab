package service

import (
	"context"
	"math"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type ServiceReliability struct {
	Service          string  `json:"service"`
	TraceCount       int     `json:"traceCount"`
	ErrorCount       int     `json:"errorCount"`
	ErrorRate        float64 `json:"errorRate"`
	AverageLatencyMs float64 `json:"averageLatencyMs"`
	P95LatencyMs     float64 `json:"p95LatencyMs"`
	HealthScore      int     `json:"healthScore"`
}

type ReliabilityService struct{ repository ports.Repository }

func NewReliabilityService(repository ports.Repository) ReliabilityService {
	return ReliabilityService{repository: repository}
}

func (s ReliabilityService) Summaries(ctx context.Context) ([]ServiceReliability, error) {
	services, err := s.repository.ListServices(ctx)
	if err != nil {
		return nil, err
	}
	traces, err := s.repository.ListTraces(ctx, ports.TraceFilter{})
	if err != nil {
		return nil, err
	}
	result := make([]ServiceReliability, 0, len(services))
	for _, item := range services {
		result = append(result, summarizeService(item.Name, traces))
	}
	return result, nil
}

func summarizeService(name string, traces []domain.Trace) ServiceReliability {
	latencies, errors := make([]float64, 0), 0
	for _, trace := range traces {
		for _, span := range trace.Spans {
			if span.Service != name {
				continue
			}
			latencies = append(latencies, span.DurationMillis())
			if span.Status == domain.SpanError {
				errors++
			}
		}
	}
	if len(latencies) == 0 {
		return ServiceReliability{Service: name, HealthScore: 100}
	}
	for i := 1; i < len(latencies); i++ {
		for j := i; j > 0 && latencies[j] < latencies[j-1]; j-- {
			latencies[j], latencies[j-1] = latencies[j-1], latencies[j]
		}
	}
	total := 0.0
	for _, value := range latencies {
		total += value
	}
	errorRate := float64(errors) / float64(len(latencies)) * 100
	p95Index := int(math.Ceil(float64(len(latencies))*.95)) - 1
	score := int(100 - errorRate*2 - latencies[p95Index]/100)
	if score < 0 {
		score = 0
	}
	return ServiceReliability{Service: name, TraceCount: len(latencies), ErrorCount: errors, ErrorRate: errorRate, AverageLatencyMs: total / float64(len(latencies)), P95LatencyMs: latencies[p95Index], HealthScore: score}
}
