package domain

import (
	"errors"
	"time"
)

type SLO struct {
	ID          string        `json:"id"`
	Service     string        `json:"service"`
	Name        string        `json:"name"`
	Target      float64       `json:"target"`
	Window      time.Duration `json:"window"`
	GoodEvents  int64         `json:"goodEvents"`
	TotalEvents int64         `json:"totalEvents"`
}

func (s SLO) Validate() error {
	if s.ID == "" || s.Service == "" || s.Name == "" {
		return errors.New("SLO identity is required")
	}
	if s.Target <= 0 || s.Target >= 100 {
		return errors.New("SLO target must be between zero and one hundred")
	}
	if s.TotalEvents < 0 || s.GoodEvents < 0 || s.GoodEvents > s.TotalEvents {
		return errors.New("SLO event counts are invalid")
	}
	return nil
}

func (s SLO) Attainment() float64 {
	if s.TotalEvents == 0 {
		return 100
	}
	return float64(s.GoodEvents) / float64(s.TotalEvents) * 100
}

func (s SLO) ErrorBudgetRemaining() float64 {
	allowed := float64(s.TotalEvents) * (1 - s.Target/100)
	if allowed == 0 {
		return 100
	}
	bad := float64(s.TotalEvents - s.GoodEvents)
	remaining := (allowed - bad) / allowed * 100
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s SLO) BurnRate() float64 {
	if s.TotalEvents == 0 {
		return 0
	}
	allowedRate := 1 - s.Target/100
	return (float64(s.TotalEvents-s.GoodEvents) / float64(s.TotalEvents)) / allowedRate
}
