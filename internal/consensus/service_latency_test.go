package consensus

import (
    "context"
    "strings"
    "testing"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/bus"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Test that publishing an event leads to a processing latency summary entry.
func TestConsensus_ProcessingLatencySummary(t *testing.T) {
    metrics.Reset()
    b := bus.New(8)
    s := NewWithSub(b.Subscribe())
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := s.Start(ctx); err != nil { t.Fatalf("start: %v", err) }
    b.Publish(ctx, bus.Event{Kind: bus.KindDuty})
    time.Sleep(20 * time.Millisecond)
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `consensus_proc_ms_count{kind="duty"}`) {
        t.Fatalf("expected consensus_proc_ms count for duty, got %q", dump)
    }
}