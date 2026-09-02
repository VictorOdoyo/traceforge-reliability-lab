package domain

import (
	"testing"
	"time"
)

func TestIncidentLifecycle(t *testing.T) {
	incident := Incident{ID: "inc-1", Status: IncidentInvestigating}
	now := time.Now()
	if err := incident.Transition(IncidentIdentified, "commander", now); err != nil {
		t.Fatal(err)
	}
	if err := incident.Transition(IncidentMonitoring, "commander", now); err != nil {
		t.Fatal(err)
	}
	if err := incident.Transition(IncidentResolved, "commander", now); err != nil {
		t.Fatal(err)
	}
	if incident.ResolvedAt == nil {
		t.Fatal("resolution time not recorded")
	}
}
func TestIncidentRejectsBackwardTransition(t *testing.T) {
	incident := Incident{ID: "inc-1", Status: IncidentMonitoring}
	if err := incident.Transition(IncidentIdentified, "commander", time.Now()); err == nil {
		t.Fatal("expected invalid transition")
	}
}
