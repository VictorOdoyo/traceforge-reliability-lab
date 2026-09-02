package httpapi

import (
	"net/http"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
)

func (s *Server) deploymentCorrelations(w http.ResponseWriter, r *http.Request) {
	items, err := s.correlations.Deployments(r.Context(), 2*time.Hour)
	if err != nil {
		writeProblem(w, 500, "Correlations unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	items, err := s.alerts.ApplySuppressions(r.Context(), s.clock())
	if err != nil {
		writeProblem(w, 500, "Alerts unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (s *Server) acknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Owner string `json:"owner"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeProblem(w, 400, "Invalid acknowledgement", err.Error())
		return
	}
	alert, err := s.alerts.Acknowledge(r.Context(), r.PathValue("id"), body.Owner, s.clock())
	if err != nil {
		writeProblem(w, 422, "Alert not acknowledged", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, alert)
}
func (s *Server) listIncidents(w http.ResponseWriter, r *http.Request) {
	items, err := s.repository.ListIncidents(r.Context())
	if err != nil {
		writeProblem(w, 500, "Incidents unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (s *Server) transitionIncident(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status domain.IncidentStatus `json:"status"`
		Actor  string                `json:"actor"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeProblem(w, 400, "Invalid transition", err.Error())
		return
	}
	incident, err := s.incidents.Transition(r.Context(), r.PathValue("id"), body.Status, body.Actor, s.clock())
	if err != nil {
		writeProblem(w, 422, "Incident not updated", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, incident)
}
