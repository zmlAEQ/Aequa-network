package monitoring

import (
    "context"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestMonitoring_ServiceOpSummary(t *testing.T) {
    metrics.Reset()
    s := New("127.0.0.1:0")
    if err := s.Start(context.Background()); err != nil {
        t.Fatalf("start: %v", err)
    }
    if err := s.Stop(context.Background()); err != nil {
        t.Fatalf("stop: %v", err)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `service_op_ms_count{op="start",service="monitoring"}`) {
        t.Fatalf("expected service_op_ms count for monitoring start, got %q", dump)
    }
    if !strings.Contains(dump, `service_op_ms_count{op="stop",service="monitoring"}`) {
        t.Fatalf("expected service_op_ms count for monitoring stop, got %q", dump)
    }
}