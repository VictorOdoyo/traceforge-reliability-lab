package service

import (
	"context"
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/platform/memory"
)

func TestAlertAcknowledgementPersistsOwner(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	repository := &memory.Repository{Alerts: []domain.Alert{{ID: "a-1", Service: "payments", State: domain.AlertFiring}}}

	alert, err := NewAlertService(repository).Acknowledge(context.Background(), "a-1", "Victor", now)
	if err != nil {
		t.Fatalf("Acknowledge() error = %v", err)
	}
	if alert.State != domain.AlertAcknowledged || alert.Owner != "Victor" {
		t.Fatalf("unexpected alert: %+v", alert)
	}
	persisted, _ := repository.ListAlerts(context.Background())
	if persisted[0].UpdatedAt != now {
		t.Fatal("acknowledgement was not persisted")
	}
}

func TestAlertSuppressionsApplyOnlyInsideWindow(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	repository := &memory.Repository{
		Alerts:       []domain.Alert{{ID: "a-1", Service: "payments", State: domain.AlertFiring}, {ID: "a-2", Service: "checkout", State: domain.AlertFiring}},
		Suppressions: []domain.SuppressionRule{{ID: "s-1", Service: "payments", StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}},
	}

	alerts, err := NewAlertService(repository).ApplySuppressions(context.Background(), now)
	if err != nil {
		t.Fatalf("ApplySuppressions() error = %v", err)
	}
	if alerts[0].State != domain.AlertSuppressed || alerts[1].State != domain.AlertFiring {
		t.Fatalf("unexpected suppression result: %+v", alerts)
	}
}

func TestAlertAcknowledgementRejectsUnknownAlert(t *testing.T) {
	_, err := NewAlertService(&memory.Repository{}).Acknowledge(context.Background(), "missing", "Victor", time.Now())
	if err == nil {
		t.Fatal("expected not found error")
	}
}
