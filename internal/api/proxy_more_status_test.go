package api

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestProxy_Upstream400Metrics(t *testing.T) {
    metrics.Reset()
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = w.Write([]byte("bad"))
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("want 400, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="400",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 400, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="proxy"}`) {
        t.Fatalf("expected latency summary for proxy, got %q", dump)
    }
}

func TestProxy_Upstream401Metrics(t *testing.T) {
    metrics.Reset()
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusUnauthorized)
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("want 401, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="401",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 401, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="proxy"}`) {
        t.Fatalf("expected latency summary for proxy, got %q", dump)
    }
}

func TestProxy_Upstream429Metrics(t *testing.T) {
    metrics.Reset()
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusTooManyRequests)
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusTooManyRequests {
        t.Fatalf("want 429, got %d", rr.Code)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="429",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 429, got %q", dump)
    }
    if !strings.Contains(dump, `api_latency_ms_count{route="proxy"}`) {
        t.Fatalf("expected latency summary for proxy, got %q", dump)
    }
}