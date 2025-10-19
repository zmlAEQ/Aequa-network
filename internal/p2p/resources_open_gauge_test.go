package p2p

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestResourceManager_OpenGaugeReflectsCurrent(t *testing.T) {
    metrics.Reset()
    r := NewResourceManager(ResourceLimits{MaxConns: 3})
    if !r.TryOpen() { t.Fatalf("first open should succeed") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, "p2p_conns_open 2") {
        t.Fatalf("want p2p_conns_open 2, got %q", dump)
    }
    r.Close(); r.Close()
    dump2 := metrics.DumpProm()
    if !strings.Contains(dump2, "p2p_conns_open 0") {
        t.Fatalf("want p2p_conns_open 0, got %q", dump2)
    }
}