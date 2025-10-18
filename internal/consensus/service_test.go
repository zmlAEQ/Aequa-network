package consensus

import (
    "context"
    "strings"
    "testing"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/bus"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestService_ConsumesEventsAndIncrementsMetrics(t *testing.T) {
    b := bus.New(4)
    s := NewWithSub(b.Subscribe())
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := s.Start(ctx); err != nil { t.Fatalf("start: %v", err) }

    // Publish one duty event
    b.Publish(ctx, bus.Event{Kind: bus.KindDuty})

    // Wait briefly for goroutine to process
    time.Sleep(50 * time.Millisecond)

    dump := metrics.DumpProm()
    if !strings.Contains(dump, "consensus_events_total{kind=\"duty\"} 1") {
        t.Fatalf("metrics not found or wrong value: %s", dump)
    }
}
