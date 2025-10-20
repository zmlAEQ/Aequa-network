package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// When a preprepare comes from a wrong leader, Process must return error,
// keep the current phase unchanged, and must not increment transitions counter.
func TestState_Preprepare_WrongLeader_DoesNotTransition(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "leader1"}
    // Start with empty phase
    if st.Phase != "" { t.Fatalf("unexpected initial phase: %q", st.Phase) }

    msg := Message{ID: "pp-bad", From: "leaderX", Type: MsgPreprepare, Height: 10, Round: 0}
    if err := st.Process(msg); err == nil {
        t.Fatalf("expected error for unauthorized leader")
    }
    if st.Phase != "" {
        t.Fatalf("phase should remain unchanged on error, got %q", st.Phase)
    }

    dump := metrics.DumpProm()
    if strings.Contains(dump, `qbft_state_transitions_total{type="preprepare"}`) {
        t.Fatalf("unexpected transition counter for preprepare on error: %q", dump)
    }
}