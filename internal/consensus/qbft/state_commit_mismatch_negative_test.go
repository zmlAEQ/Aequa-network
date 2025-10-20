package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Commit with a mismatching proposal ID must return error, keep phase, and
// must not increase commit transition metrics.
func TestState_Commit_MismatchProposal_DoesNotCountOrTransition(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Establish proposal id "blk1" and reach prepared.
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 10, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 10, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 10, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" {
        t.Fatalf("want prepared, got %q", st.Phase)
    }

    // Commit for a different proposal id should fail and not change phase.
    if err := st.Process(Message{ID: "blkX", From: "C1", Type: MsgCommit, Height: 10, Round: 1}); err == nil {
        t.Fatalf("expected error for mismatched proposal id on commit")
    }
    if st.Phase != "prepared" {
        t.Fatalf("phase changed unexpectedly: %q", st.Phase)
    }

    dump := metrics.DumpProm()
    if strings.Contains(dump, `qbft_state_transitions_total{type="commit"}`) {
        t.Fatalf("unexpected commit transition count on error: %q", dump)
    }
}

