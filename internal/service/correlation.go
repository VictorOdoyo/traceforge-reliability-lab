package service

import (
	"context"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type DeploymentCorrelation struct {
	Deployment  domain.Deployment `json:"deployment"`
	ErrorTraces int               `json:"errorTraces"`
	TotalTraces int               `json:"totalTraces"`
	ErrorRate   float64           `json:"errorRate"`
}
type CorrelationService struct{ repository ports.Repository }

func NewCorrelationService(repository ports.Repository) CorrelationService {
	return CorrelationService{repository: repository}
}
func (s CorrelationService) Deployments(ctx context.Context, window time.Duration) ([]DeploymentCorrelation, error) {
	deployments, err := s.repository.ListDeployments(ctx)
	if err != nil {
		return nil, err
	}
	traces, err := s.repository.ListTraces(ctx, ports.TraceFilter{})
	if err != nil {
		return nil, err
	}
	result := make([]DeploymentCorrelation, 0, len(deployments))
	for _, deployment := range deployments {
		item := DeploymentCorrelation{Deployment: deployment}
		for _, trace := range traces {
			if trace.RootService == deployment.Service && deployment.Correlates(trace.StartedAt, window) {
				item.TotalTraces++
				if trace.HasError {
					item.ErrorTraces++
				}
			}
		}
		if item.TotalTraces > 0 {
			item.ErrorRate = float64(item.ErrorTraces) / float64(item.TotalTraces) * 100
		}
		result = append(result, item)
	}
	return result, nil
}
