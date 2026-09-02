package httpapi

import (
	"net/http"
	"strconv"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
)

func (s *Server) searchTraces(w http.ResponseWriter, r *http.Request) {
	minimum, _ := strconv.ParseFloat(r.URL.Query().Get("minDurationMs"), 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 50
	}
	if limit < 1 || limit > 200 {
		writeProblem(w, 400, "Invalid limit", "Limit must be between 1 and 200.")
		return
	}
	traces, err := s.traces.Search(r.Context(), r.URL.Query().Get("service"), r.URL.Query().Get("status"), minimum, limit)
	if err != nil {
		writeProblem(w, 500, "Traces unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, traces)
}
func (s *Server) ingestTrace(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID    string        `json:"id"`
		Spans []domain.Span `json:"spans"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeProblem(w, 400, "Invalid trace", err.Error())
		return
	}
	trace, err := domain.BuildTrace(body.ID, body.Spans)
	if err != nil {
		writeProblem(w, 422, "Trace rejected", err.Error())
		return
	}
	if err := s.repository.SaveTrace(r.Context(), trace); err != nil {
		writeProblem(w, 500, "Trace storage failed", err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, trace)
}
func (s *Server) serviceTopology(w http.ResponseWriter, r *http.Request) {
	edges, err := s.topology.Edges(r.Context())
	if err != nil {
		writeProblem(w, 500, "Topology unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, edges)
}
