package service

import (
	"context"
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/platform/memory"
)

func TestIncidentTransitionAddsTimelineAndPersists(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	repository := &memory.Repository{Incidents: []domain.Incident{{ID: "inc-1", Status: domain.IncidentInvestigating}}}

	incident, err := NewIncidentService(repository).Transition(context.Background(), "inc-1", domain.IncidentIdentified, "Samira", now)
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}
	if incident.Status != domain.IncidentIdentified || len(incident.Timeline) != 1 || incident.Timeline[0].Actor != "Samira" {
		t.Fatalf("unexpected incident: %+v", incident)
	}
	persisted, _ := repository.ListIncidents(context.Background())
	if persisted[0].Status != domain.IncidentIdentified {
		t.Fatal("transition was not persisted")
	}
}

func TestResolvedTransitionClosesWithoutTimelineAppend(t *testing.T) {
	repository := &memory.Repository{Incidents: []domain.Incident{{ID: "inc-1", Status: domain.IncidentMonitoring}}}
	incident, err := NewIncidentService(repository).Transition(context.Background(), "inc-1", domain.IncidentResolved, "Victor", time.Now())
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}
	if incident.ResolvedAt == nil || len(incident.Timeline) != 0 {
		t.Fatalf("unexpected resolved incident: %+v", incident)
	}
}

func TestIncidentTransitionRejectsUnknownIncident(t *testing.T) {
	_, err := NewIncidentService(&memory.Repository{}).Transition(context.Background(), "missing", domain.IncidentResolved, "Victor", time.Now())
	if err == nil {
		t.Fatal("expected not found error")
	}
}
