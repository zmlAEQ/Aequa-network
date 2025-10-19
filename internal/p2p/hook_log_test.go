package p2p

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestLogHook_IncrementsCounters(t *testing.T) {
    metrics.Reset()
    h := LogHook{}
    h.OnPeerJoin("A")
    h.OnPeerLeave("A")
    dump := metrics.DumpProm()
    if !strings.Contains(dump, "p2p_peers_joined_total 1") {
        t.Fatalf("want joined=1, got %q", dump)
    }
    if !strings.Contains(dump, "p2p_peers_left_total 1") {
        t.Fatalf("want left=1, got %q", dump)
    }
}

func TestService_WithLogHook_RecordsJoinLeave(t *testing.T) {
    metrics.Reset()
    s := NewWithOpts(nil, AllowAllGate{}, NewResourceManager(ResourceLimits{MaxConns: 2}), LogHook{})
    if err := s.Connect("A"); err != nil { t.Fatalf("connect: %v", err) }
    s.Disconnect("A")
    dump := metrics.DumpProm()
    if !strings.Contains(dump, "p2p_peers_joined_total 1") || !strings.Contains(dump, "p2p_peers_left_total 1") {
        t.Fatalf("want joined=1 & left=1, got %q", dump)
    }
}