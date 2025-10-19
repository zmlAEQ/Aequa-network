package p2p

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestAllowListGate_ServiceAttempts(t *testing.T) {
    metrics.Reset()
    g := NewAllowListGate("A")
    s := NewWithOpts(nil, g, NewResourceManager(ResourceLimits{MaxConns: 2}), NopHook{})

    if err := s.Connect("A"); err != nil { t.Fatalf("A should pass: %v", err) }
    if err := s.Connect("B"); err == nil { t.Fatalf("B should be denied") }

    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="allowed"} 1`) {
        t.Fatalf("expected allowed=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="denied"} 1`) {
        t.Fatalf("expected denied=1, got %q", dump)
    }
}