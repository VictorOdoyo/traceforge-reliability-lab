package domain

import (
	"errors"
	"time"
)

type IncidentStatus string

const (
	IncidentInvestigating IncidentStatus = "investigating"
	IncidentIdentified    IncidentStatus = "identified"
	IncidentMonitoring    IncidentStatus = "monitoring"
	IncidentResolved      IncidentStatus = "resolved"
)

type TimelineEvent struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	Message    string    `json:"message"`
	Actor      string    `json:"actor"`
	OccurredAt time.Time `json:"occurredAt"`
}

type Incident struct {
	ID               string          `json:"id"`
	Title            string          `json:"title"`
	Severity         Severity        `json:"severity"`
	Status           IncidentStatus  `json:"status"`
	Commander        string          `json:"commander"`
	AffectedServices []string        `json:"affectedServices"`
	StartedAt        time.Time       `json:"startedAt"`
	ResolvedAt       *time.Time      `json:"resolvedAt,omitempty"`
	Timeline         []TimelineEvent `json:"timeline"`
}

func (i *Incident) AddEvent(event TimelineEvent) error {
	if event.ID == "" || event.Message == "" || event.Actor == "" {
		return errors.New("timeline event identity, message, and actor are required")
	}
	if i.Status == IncidentResolved {
		return errors.New("resolved incident timeline is closed")
	}
	i.Timeline = append(i.Timeline, event)
	return nil
}

func (i *Incident) Transition(status IncidentStatus, actor string, at time.Time) error {
	if actor == "" {
		return errors.New("transition actor is required")
	}
	allowed := map[IncidentStatus][]IncidentStatus{IncidentInvestigating: {IncidentIdentified, IncidentMonitoring, IncidentResolved}, IncidentIdentified: {IncidentMonitoring, IncidentResolved}, IncidentMonitoring: {IncidentResolved}}
	valid := false
	for _, next := range allowed[i.Status] {
		if next == status {
			valid = true
		}
	}
	if !valid {
		return errors.New("invalid incident transition")
	}
	i.Status = status
	if status == IncidentResolved {
		i.ResolvedAt = &at
	}
	return nil
}
