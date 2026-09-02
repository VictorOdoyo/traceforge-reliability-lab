package httpapi

import (
	"net/http"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
)

type Overview struct {
	ServiceCount     int `json:"serviceCount"`
	OpenAlerts       int `json:"openAlerts"`
	ActiveIncidents  int `json:"activeIncidents"`
	ErrorBudgetRisks int `json:"errorBudgetRisks"`
	Services         any `json:"services"`
}

func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	summaries, err := s.reliability.Summaries(r.Context())
	if err != nil {
		writeProblem(w, 500, "Overview unavailable", err.Error())
		return
	}
	alerts, _ := s.repository.ListAlerts(r.Context())
	incidents, _ := s.repository.ListIncidents(r.Context())
	slos, _ := s.repository.ListSLOs(r.Context())
	result := Overview{ServiceCount: len(summaries), Services: summaries}
	for _, alert := range alerts {
		if alert.State != domain.AlertResolved && alert.State != domain.AlertSuppressed {
			result.OpenAlerts++
		}
	}
	for _, incident := range incidents {
		if incident.Status != domain.IncidentResolved {
			result.ActiveIncidents++
		}
	}
	for _, slo := range slos {
		if slo.BurnRate() > 1 {
			result.ErrorBudgetRisks++
		}
	}
	writeJSON(w, http.StatusOK, result)
}
func (s *Server) services(w http.ResponseWriter, r *http.Request) {
	services, err := s.repository.ListServices(r.Context())
	if err != nil {
		writeProblem(w, 500, "Services unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, services)
}
func (s *Server) slos(w http.ResponseWriter, r *http.Request) {
	slos, err := s.repository.ListSLOs(r.Context())
	if err != nil {
		writeProblem(w, 500, "SLOs unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, slos)
}
