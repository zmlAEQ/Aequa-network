package api

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Ensure 4xx from upstream is surfaced and metrics are recorded with the 4xx code.
func TestProxy_UpstreamClientErrorMetrics(t *testing.T) {
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.NotFound(w, r)
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusNotFound {
        t.Fatalf("want 404, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="404",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 404, got %q", dump)
    }
}