package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// Duplicate commit must be a no-op (phase unchanged) but still be counted as a
// processed message in qbft_state_transitions_total{type="commit"}.
func TestState_Commit_Duplicate_NoOpCounts(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Reach prepared
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 6, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 6, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 6, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("want prepared, got %q", st.Phase) }

    // First commit advances to commit.
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 6, Round: 1}); err != nil {
        t.Fatalf("commit1: %v", err)
    }
    if st.Phase != "commit" { t.Fatalf("want commit, got %q", st.Phase) }

    // Duplicate commit is a no-op but still counts as processed.
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 6, Round: 1}); err != nil {
        t.Fatalf("duplicate commit must not error: %v", err)
    }
    if st.Phase != "commit" { t.Fatalf("duplicate must not change phase: %q", st.Phase) }

    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_state_transitions_total{type="commit"} 2`) {
        t.Fatalf("duplicate commit should count as processed (want 2), got %q", dump)
    }
}

// Commit with mismatched proposal ID must return error, keep phase and must not
// increment commit transitions counter (only successful/non-error paths count).
func TestState_Commit_Mismatch_NoCount(t *testing.T) {
    metrics.Reset()
    st := &State{Leader: "L"}

    // Reach prepared
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 7, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 7, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 7, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("want prepared, got %q", st.Phase) }

    // First commit (ok)
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 7, Round: 1}); err != nil {
        t.Fatalf("commit1: %v", err)
    }
    // Mismatched commit (error)
    if err := st.Process(Message{ID: "blkX", From: "C2", Type: MsgCommit, Height: 7, Round: 1}); err == nil {
        t.Fatalf("expected error for mismatched proposal id on commit")
    }
    if st.Phase != "commit" { t.Fatalf("phase changed unexpectedly: %q", st.Phase) }

    dump := metrics.DumpProm()
    if strings.Contains(dump, `qbft_state_transitions_total{type="commit"} 2`) {
        t.Fatalf("error path must not increment commit counter: %q", dump)
    }
}

