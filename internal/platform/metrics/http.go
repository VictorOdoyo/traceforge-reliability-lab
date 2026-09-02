package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTPMetrics struct {
	registry *prometheus.Registry
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
	inflight prometheus.Gauge
}

func NewHTTPMetrics() *HTTPMetrics {
	registry := prometheus.NewRegistry()
	requests := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "traceforge_http_requests_total", Help: "HTTP requests handled by status and method."}, []string{"method", "status"})
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "traceforge_http_request_duration_seconds", Help: "HTTP request latency in seconds.", Buckets: prometheus.DefBuckets}, []string{"method"})
	inflight := prometheus.NewGauge(prometheus.GaugeOpts{Name: "traceforge_http_inflight_requests", Help: "Requests currently being handled."})
	registry.MustRegister(requests, duration, inflight, prometheus.NewGoCollector())
	return &HTTPMetrics{registry: registry, requests: requests, duration: duration, inflight: inflight}
}
func (m *HTTPMetrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
func (m *HTTPMetrics) Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		m.inflight.Inc()
		defer m.inflight.Dec()
		recorder := &metricRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		m.requests.WithLabelValues(r.Method, strconv.Itoa(recorder.status)).Inc()
		m.duration.WithLabelValues(r.Method).Observe(time.Since(started).Seconds())
	})
}

type metricRecorder struct {
	http.ResponseWriter
	status int
}

func (r *metricRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
