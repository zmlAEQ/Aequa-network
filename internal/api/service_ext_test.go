package api

import (
    "bytes"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestHandleDuty_InvalidJSON(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader([]byte("{")))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("want 400, got %d", rr.Code)
    }
}

func TestHandleDuty_InvalidType(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    body := []byte(`{"type":"x","height":1,"round":0,"payload":{}}`)
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(body))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("want 400, got %d", rr.Code)
    }
}

func TestHandleDuty_RoundOutOfRange(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    bigRound := uint64(1) << 41
    body := []byte(fmt.Sprintf(`{"type":"attester","height":1,"round":%d,"payload":{}}`, bigRound))
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(body))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("want 400, got %d", rr.Code)
    }
}

func TestProxy_NotFound_NoUpstream(t *testing.T) {
    s := &Service{upstream: ""}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    if rr.Code != http.StatusNotFound {
        t.Fatalf("want 404, got %d", rr.Code)
    }
}

func TestMetrics_DutyAccepted_Increments(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    body := []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`)
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(body))
    s.handleDuty(rr, req)
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="202",route="/v1/duty"}`) {
        t.Fatalf("expected metrics increment for /v1/duty 202, got %q", dump)
    }
}

func TestMetrics_ProxySuccess_Increments(t *testing.T) {
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("ok"))
    }))
    defer backend.Close()

    s := &Service{upstream: backend.URL}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="200",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 200, got %q", dump)
    }
}

func TestMetrics_ProxyError_Increments(t *testing.T) {
    s := &Service{upstream: "http://127.0.0.1:1"}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    s.proxy(rr, req)
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `api_requests_total{code="502",route="proxy"}`) {
        t.Fatalf("expected metrics increment for proxy 502, got %q", dump)
    }
}