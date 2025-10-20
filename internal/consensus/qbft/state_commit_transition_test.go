package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Minimal placeholder: after entering prepared, a commit message should
// advance phase to commit and emit one transition counter for commit.
func TestState_Commit_AfterPrepared(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Preprepare establishes proposal id.
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 6, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("want preprepared, got %q", st.Phase)
    }

    // Two prepares reach threshold -> prepared (as per simplified model).
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 6, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("still preprepared after first prepare, got %q", st.Phase)
    }
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 6, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" {
        t.Fatalf("want prepared, got %q", st.Phase)
    }

    // Commit message should advance to commit and record metrics.
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 6, Round: 1}); err != nil {
        t.Fatalf("commit: %v", err)
    }
    if st.Phase != "commit" {
        t.Fatalf("want commit, got %q", st.Phase)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_state_transitions_total{type="commit"} 1`) {
        t.Fatalf("missing commit transition metric: %q", dump)
    }
}