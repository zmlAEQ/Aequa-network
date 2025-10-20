package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Commit received before reaching prepared must not advance the phase and
// must not count a commit transition. It should return an error in a strict
// implementation; until the logic is implemented, this test documents the
// intended behavior as a guardrail.
func TestState_Commit_BeforePrepared_DoesNotAdvanceOrCount(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Establish proposal id but stay in preprepared (no prepare threshold yet).
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 9, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("want preprepared, got %q", st.Phase)
    }

    // Receiving commit before prepared must NOT advance to commit, and must
    // NOT increase the commit transition metric.
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 9, Round: 1}); err == nil {
        // In a strict implementation this should be an error; allow the guard
        // to enforce no silent success.
        t.Fatalf("expected error for commit before prepared")
    }
    if st.Phase != "preprepared" {
        t.Fatalf("phase changed unexpectedly: %q", st.Phase)
    }
    dump := metrics.DumpProm()
    if strings.Contains(dump, `qbft_state_transitions_total{type="commit"}`) {
        t.Fatalf("unexpected commit transition count on error: %q", dump)
    }
}