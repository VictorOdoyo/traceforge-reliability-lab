package service

import (
	"context"
	"errors"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type IncidentService struct{ repository ports.Repository }

func NewIncidentService(repository ports.Repository) IncidentService {
	return IncidentService{repository: repository}
}
func (s IncidentService) Transition(ctx context.Context, id string, status domain.IncidentStatus, actor string, at time.Time) (domain.Incident, error) {
	incidents, err := s.repository.ListIncidents(ctx)
	if err != nil {
		return domain.Incident{}, err
	}
	for _, incident := range incidents {
		if incident.ID != id {
			continue
		}
		if err := incident.Transition(status, actor, at); err != nil {
			return domain.Incident{}, err
		}
		event := domain.TimelineEvent{ID: id + "-" + string(status), Kind: "status", Message: "Incident moved to " + string(status), Actor: actor, OccurredAt: at}
		if status != domain.IncidentResolved {
			_ = incident.AddEvent(event)
		}
		if err := s.repository.SaveIncident(ctx, incident); err != nil {
			return domain.Incident{}, err
		}
		return incident, nil
	}
	return domain.Incident{}, errors.New("incident not found")
}
