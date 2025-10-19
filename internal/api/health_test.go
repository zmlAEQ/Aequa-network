package api

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestHandleHealth_LogsAndSummary(t *testing.T) {
    metrics.Reset()
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    s.handleHealth(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("want 200, got %d", rr.Code)
    }
    if body := rr.Body.String(); body != "ok" {
        t.Fatalf("want body 'ok', got %q", body)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="200",route="/health"}`) {
        t.Fatalf("expected metrics increment for /health 200, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="/health"}`) {
        t.Fatalf("expected latency summary for /health, got %q", dump)
    }
}