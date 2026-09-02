package service

import (
	"context"
	"errors"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type AlertService struct{ repository ports.Repository }

func NewAlertService(repository ports.Repository) AlertService {
	return AlertService{repository: repository}
}
func (s AlertService) Acknowledge(ctx context.Context, id, owner string, at time.Time) (domain.Alert, error) {
	alerts, err := s.repository.ListAlerts(ctx)
	if err != nil {
		return domain.Alert{}, err
	}
	for _, alert := range alerts {
		if alert.ID != id {
			continue
		}
		if err := alert.Acknowledge(owner, at); err != nil {
			return domain.Alert{}, err
		}
		if err := s.repository.SaveAlert(ctx, alert); err != nil {
			return domain.Alert{}, err
		}
		return alert, nil
	}
	return domain.Alert{}, errors.New("alert not found")
}
func (s AlertService) ApplySuppressions(ctx context.Context, at time.Time) ([]domain.Alert, error) {
	alerts, err := s.repository.ListAlerts(ctx)
	if err != nil {
		return nil, err
	}
	rules, err := s.repository.ListSuppressions(ctx)
	if err != nil {
		return nil, err
	}
	for index := range alerts {
		for _, rule := range rules {
			if rule.Matches(alerts[index], at) {
				alerts[index].State = domain.AlertSuppressed
				alerts[index].UpdatedAt = at
				_ = s.repository.SaveAlert(ctx, alerts[index])
				break
			}
		}
	}
	return alerts, nil
}
