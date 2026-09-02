package domain

import (
	"errors"
	"time"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

type AlertState string

const (
	AlertFiring       AlertState = "firing"
	AlertAcknowledged AlertState = "acknowledged"
	AlertResolved     AlertState = "resolved"
	AlertSuppressed   AlertState = "suppressed"
)

type Alert struct {
	ID        string            `json:"id"`
	Service   string            `json:"service"`
	Title     string            `json:"title"`
	Severity  Severity          `json:"severity"`
	State     AlertState        `json:"state"`
	StartedAt time.Time         `json:"startedAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Owner     string            `json:"owner,omitempty"`
	Labels    map[string]string `json:"labels"`
}

func (a *Alert) Acknowledge(owner string, at time.Time) error {
	if a.State != AlertFiring {
		return errors.New("only firing alerts can be acknowledged")
	}
	if owner == "" {
		return errors.New("alert owner is required")
	}
	a.State, a.Owner, a.UpdatedAt = AlertAcknowledged, owner, at
	return nil
}

func (a *Alert) Resolve(at time.Time) error {
	if a.State == AlertResolved {
		return errors.New("alert is already resolved")
	}
	a.State, a.UpdatedAt = AlertResolved, at
	return nil
}
