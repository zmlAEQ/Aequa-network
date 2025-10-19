package p2p

import (
    "testing"
    "strings"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestManager_PeersCounters(t *testing.T) {
    metrics.Reset()
    m := NewManager()
    m.AddPeer("A")
    m.AddPeer("B")
    m.RemovePeer("A")
    dump := metrics.DumpProm()
    if !strings.Contains(dump, "p2p_peers_added_total 2") {
        t.Fatalf("want added=2, got %q", dump)
    }
    if !strings.Contains(dump, "p2p_peers_removed_total 1") {
        t.Fatalf("want removed=1, got %q", dump)
    }
}

func TestManager_BroadcastMetrics(t *testing.T) {
    metrics.Reset()
    m := NewManager()
    m.AddPeer("A"); m.AddPeer("B")
    n := m.Broadcast(Message{From: "A", Payload: []byte("x"), TraceID: "tid"})
    if n != 2 { t.Fatalf("want 2 peers, got %d", n) }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_messages_total{kind="broadcast"} 1`) {
        t.Fatalf("expected messages_total inc, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_broadcast_ms_count{kind="broadcast"}`) {
        t.Fatalf("expected broadcast summary, got %q", dump)
    }
}