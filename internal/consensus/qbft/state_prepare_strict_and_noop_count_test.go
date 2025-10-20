package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Prepare without a prior preprepare must be rejected, phase unchanged, and
// must not increment the prepare transition counter.
func TestState_Prepare_WithoutPreprepare_ErrorNoCount(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 1, Round: 1}); err == nil {
        t.Fatalf("expected error for prepare before preprepare")
    }
    if st.Phase != "" { // phase should remain default (no transition)
        t.Fatalf("phase changed unexpectedly: %q", st.Phase)
    }

    dump := metrics.DumpProm()
    if strings.Contains(dump, `qbft_state_transitions_total{type="prepare"}`) {
        t.Fatalf("unexpected prepare transition count on error: %q", dump)
    }
}

// Duplicate prepare (same From, same proposal) must be a no-op but still count
// as a processed message (non-error) in the transitions counter, while keeping
// the phase unchanged below threshold.
func TestState_Prepare_Duplicate_NoOpCounts(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Establish proposal context.
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 2, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("want preprepared, got %q", st.Phase)
    }

    // First prepare from P1 increments prepare counter but does not reach threshold.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 2, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("still preprepared after first prepare, got %q", st.Phase)
    }

    // Duplicate from P1 is a no-op; should not advance phase, but still count.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 2, Round: 1}); err != nil {
        t.Fatalf("duplicate prepare must not error: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("duplicate prepare must not change phase, got %q", st.Phase)
    }

    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_state_transitions_total{type="prepare"} 2`) {
        t.Fatalf("duplicate prepare should count as processed (want 2), got %q", dump)
    }
}

