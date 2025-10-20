package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Prepare mismatched with current ProposedID must return error, keep phase, and not increase transitions.
func TestState_Prepare_MismatchProposal_DoesNotCountOrTransition(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Move to preprepared with a proposal id "blk1".
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 8, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("want preprepared, got %q", st.Phase) }

    // Prepare for a different proposal id should fail and not change phase.
    if err := st.Process(Message{ID: "blkX", From: "P1", Type: MsgPrepare, Height: 8, Round: 1}); err == nil {
        t.Fatalf("expected error for mismatched proposal id")
    }
    if st.Phase != "preprepared" { t.Fatalf("phase changed unexpectedly: %q", st.Phase) }

    dump := metrics.DumpProm()
    if strings.Contains(dump, `qbft_state_transitions_total{type="prepare"}`) {
        t.Fatalf("unexpected prepare transition count on error: %q", dump)
    }
}