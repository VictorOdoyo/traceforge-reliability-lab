package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPInstrumentationExportsRequestMetrics(t *testing.T) {
	metrics := NewHTTPMetrics()
	handler := metrics.Instrument(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/traces", nil))

	recorder := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body, _ := io.ReadAll(recorder.Result().Body)
	text := string(body)
	if !strings.Contains(text, `traceforge_http_requests_total{method="POST",status="202"} 1`) {
		t.Fatalf("request metric not exported:\n%s", text)
	}
	if !strings.Contains(text, "traceforge_http_request_duration_seconds_count") {
		t.Fatal("duration metric not exported")
	}
}
