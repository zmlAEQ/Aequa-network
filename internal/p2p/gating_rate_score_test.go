package p2p

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestRateLimitGate_AttemptsLabels(t *testing.T) {
    metrics.Reset()
    g := NewRateLimitGate(1)
    s := NewWithOpts(nil, &g, NewResourceManager(ResourceLimits{MaxConns: 3}), NopHook{})
    if err := s.Connect("A"); err != nil { t.Fatalf("A should pass: %v", err) }
    if err := s.Connect("B"); err == nil { t.Fatalf("B should be rate limited") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="allowed"} 1`) {
        t.Fatalf("expected allowed=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="rate_limited"} 1`) {
        t.Fatalf("expected rate_limited=1, got %q", dump)
    }
}

func TestScoreGate_ScoredOut(t *testing.T) {
    metrics.Reset()
    scores := map[PeerID]int64{"A":12}
    g := NewScoreGate(10, scores)
    s := NewWithOpts(nil, g, NewResourceManager(ResourceLimits{MaxConns: 2}), NopHook{})
    if err := s.Connect("A"); err != nil { t.Fatalf("A should pass: %v", err) }
    if err := s.Connect("B"); err == nil { t.Fatalf("B should be scored out") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="allowed"} 1`) {
        t.Fatalf("expected allowed=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="scored_out"} 1`) {
        t.Fatalf("expected scored_out=1, got %q", dump)
    }
}