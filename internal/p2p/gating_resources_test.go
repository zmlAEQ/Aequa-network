package p2p

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

type denyGate struct{}
func (denyGate) Allow(id PeerID) bool { return false }

func TestService_ConnectDeniedByGate(t *testing.T) {
    metrics.Reset()
    s := NewWithOpts(nil, denyGate{}, NewResourceManager(DefaultResourceLimits()), NopHook{})
    if err := s.Connect("A"); err == nil {
        t.Fatalf("expected deny error")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="denied"} 1`) {
        t.Fatalf("expected denied=1, got %q", dump)
    }
}

func TestService_ConnectLimitedByResource(t *testing.T) {
    metrics.Reset()
    r := NewResourceManager(ResourceLimits{MaxConns: 1})
    s := NewWithOpts(nil, AllowAllGate{}, r, NopHook{})
    if err := s.Connect("A"); err != nil { t.Fatalf("A should pass: %v", err) }
    if err := s.Connect("B"); err == nil { t.Fatalf("B should be limited") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="allowed"} 1`) {
        t.Fatalf("expected allowed=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="limited"} 1`) {
        t.Fatalf("expected limited=1, got %q", dump)
    }
    s.Disconnect("A")
}