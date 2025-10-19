package api

import (
    "bytes"
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/trace"
)

func TestHandleDuty_PropagatesTraceID(t *testing.T) {
    var got string
    s := &Service{onPublish: func(ctx context.Context, _ []byte) error {
        if id, ok := trace.FromContext(ctx); ok { got = id }
        return nil
    }}

    rr := httptest.NewRecorder()
    body := []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`)
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(body))
    req.Header.Set("X-Trace-ID", "abc-123")

    s.handleDuty(rr, req)
    if rr.Code != http.StatusAccepted { t.Fatalf("want 202, got %d", rr.Code) }
    if got != "abc-123" { t.Fatalf("trace id not propagated, got=%q", got) }
}