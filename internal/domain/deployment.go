package domain

import "time"

type Deployment struct {
	ID            string    `json:"id"`
	Service       string    `json:"service"`
	Version       string    `json:"version"`
	Environment   string    `json:"environment"`
	CommitSHA     string    `json:"commitSha"`
	DeployedBy    string    `json:"deployedBy"`
	DeployedAt    time.Time `json:"deployedAt"`
	ChangeSummary string    `json:"changeSummary"`
}

func (d Deployment) Correlates(eventTime time.Time, window time.Duration) bool {
	if eventTime.Before(d.DeployedAt) {
		return false
	}
	return eventTime.Sub(d.DeployedAt) <= window
}
