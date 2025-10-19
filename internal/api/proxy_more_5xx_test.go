package api

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestProxy_Upstream503Metrics(t *testing.T) {
    metrics.Reset()
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusServiceUnavailable)
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusServiceUnavailable {
        t.Fatalf("want 503, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="503",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 503, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="proxy"}`) {
        t.Fatalf("expected latency summary for proxy, got %q", dump)
    }
}

func TestProxy_Upstream504Metrics(t *testing.T) {
    metrics.Reset()
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusGatewayTimeout)
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusGatewayTimeout {
        t.Fatalf("want 504, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="504",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 504, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="proxy"}`) {
        t.Fatalf("expected latency summary for proxy, got %q", dump)
    }
}