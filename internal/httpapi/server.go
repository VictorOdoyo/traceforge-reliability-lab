package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	telemetrymetrics "github.com/VictorOdoyo/traceforge-reliability-lab/internal/platform/metrics"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/service"
)

type Server struct {
	repository   ports.Repository
	traces       service.TraceQuery
	reliability  service.ReliabilityService
	topology     service.TopologyService
	correlations service.CorrelationService
	alerts       service.AlertService
	incidents    service.IncidentService
	logger       *slog.Logger
	clock        func() time.Time
}

func New(repository ports.Repository, logger *slog.Logger) http.Handler {
	server := &Server{repository: repository, traces: service.NewTraceQuery(repository), reliability: service.NewReliabilityService(repository), topology: service.NewTopologyService(repository), correlations: service.NewCorrelationService(repository), alerts: service.NewAlertService(repository), incidents: service.NewIncidentService(repository), logger: logger, clock: time.Now}
	mux := http.NewServeMux()
	metrics := telemetrymetrics.NewHTTPMetrics()
	mux.Handle("GET /metrics", metrics.Handler())
	mux.HandleFunc("GET /health/live", server.live)
	mux.HandleFunc("GET /health/ready", server.ready)
	mux.HandleFunc("GET /api/v1/overview", server.overview)
	mux.HandleFunc("GET /api/v1/services", server.services)
	mux.HandleFunc("GET /api/v1/traces", server.searchTraces)
	mux.HandleFunc("POST /api/v1/traces", server.ingestTrace)
	mux.HandleFunc("GET /api/v1/topology", server.serviceTopology)
	mux.HandleFunc("GET /api/v1/slos", server.slos)
	mux.HandleFunc("GET /api/v1/deployments/correlations", server.deploymentCorrelations)
	mux.HandleFunc("GET /api/v1/alerts", server.listAlerts)
	mux.HandleFunc("POST /api/v1/alerts/{id}/acknowledge", server.acknowledgeAlert)
	mux.HandleFunc("GET /api/v1/incidents", server.listIncidents)
	mux.HandleFunc("POST /api/v1/incidents/{id}/transition", server.transitionIncident)
	return metrics.Instrument(requestLog(logger, recoverPanic(logger, securityHeaders(cors(mux)))))
}
