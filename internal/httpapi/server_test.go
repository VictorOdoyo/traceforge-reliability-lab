package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/seed"
)

func TestHealthAndOverviewContracts(t *testing.T) {
	handler := New(seed.Repository(time.Now()), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, path := range []string{"/health/live", "/api/v1/overview", "/api/v1/topology", "/api/v1/traces?limit=5"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("%s returned %d: %s", path, response.Code, response.Body.String())
		}
	}
}
func TestAcknowledgeAlertContract(t *testing.T) {
	handler := New(seed.Repository(time.Now()), slog.New(slog.NewTextHandler(io.Discard, nil)))
	body, _ := json.Marshal(map[string]string{"owner": "Victor"})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/alerts/alert-301/acknowledge", bytes.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["state"] != "acknowledged" {
		t.Fatalf("unexpected response: %#v", result)
	}
}
func TestTraceLimitValidation(t *testing.T) {
	handler := New(seed.Repository(time.Now()), slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/traces?limit=500", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
}
