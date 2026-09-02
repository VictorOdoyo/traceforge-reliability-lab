package domain

import "time"

type SuppressionRule struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Service    string    `json:"service,omitempty"`
	LabelKey   string    `json:"labelKey,omitempty"`
	LabelValue string    `json:"labelValue,omitempty"`
	StartsAt   time.Time `json:"startsAt"`
	EndsAt     time.Time `json:"endsAt"`
	Reason     string    `json:"reason"`
}

func (r SuppressionRule) Matches(alert Alert, at time.Time) bool {
	if at.Before(r.StartsAt) || !at.Before(r.EndsAt) {
		return false
	}
	if r.Service != "" && r.Service != alert.Service {
		return false
	}
	if r.LabelKey != "" && alert.Labels[r.LabelKey] != r.LabelValue {
		return false
	}
	return true
}
