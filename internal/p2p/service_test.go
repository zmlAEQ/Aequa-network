package p2p

import (
    "context"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestP2P_ServiceOpSummary(t *testing.T) {
    s := &Service{}
    if err := s.Start(context.Background()); err != nil {
        t.Fatalf("start: %v", err)
    }
    if err := s.Stop(context.Background()); err != nil {
        t.Fatalf("stop: %v", err)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `service_op_ms_count{op="start",service="p2p"}`) {
        t.Fatalf("expected service_op_ms count for p2p start, got %q", dump)
    }
    if !strings.Contains(dump, `service_op_ms_count{op="stop",service="p2p"}`) {
        t.Fatalf("expected service_op_ms count for p2p stop, got %q", dump)
    }
}

