package monitoring

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestHandleMetrics_LogsAndSummary(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
    s.handleMetrics(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("want 200, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="200",route="/metrics"}`) {
        t.Fatalf("expected metrics increment for /metrics 200, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="/metrics"}`) {
        t.Fatalf("expected summary count for /metrics, got %q", dump)
    }
}

