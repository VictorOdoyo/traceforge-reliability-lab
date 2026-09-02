package domain

import (
	"errors"
	"strings"
)

type Service struct {
	Name        string            `json:"name"`
	Team        string            `json:"team"`
	Tier        int               `json:"tier"`
	Environment string            `json:"environment"`
	Language    string            `json:"language"`
	Labels      map[string]string `json:"labels"`
}

func NewService(name, team, environment string, tier int) (Service, error) {
	name, team, environment = strings.TrimSpace(name), strings.TrimSpace(team), strings.TrimSpace(environment)
	if name == "" || team == "" || environment == "" {
		return Service{}, errors.New("service name, team, and environment are required")
	}
	if tier < 1 || tier > 4 {
		return Service{}, errors.New("service tier must be between 1 and 4")
	}
	return Service{Name: name, Team: team, Environment: environment, Tier: tier, Labels: map[string]string{}}, nil
}

func (s Service) Key() string { return s.Environment + "/" + s.Name }
