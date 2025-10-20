package qbft

import "testing"

// Duplicate Commit from the same From for the same proposal ID must not
// advance further nor error; it's a no-op. This test documents the dedup
// requirement as a guardrail (implementation may be added later).
func TestState_Commit_Dedup_SameFrom_NoAdvanceNoError(t *testing.T) {
    st := &State{Leader: "L"}

    // Establish proposal id and reach prepared via two prepares.
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 7, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("want preprepared, got %q", st.Phase)
    }
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 7, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("still preprepared after first prepare, got %q", st.Phase)
    }
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 7, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" {
        t.Fatalf("want prepared, got %q", st.Phase)
    }

    // First commit advances to commit.
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 7, Round: 1}); err != nil {
        t.Fatalf("commit1: %v", err)
    }
    if st.Phase != "commit" {
        t.Fatalf("want commit, got %q", st.Phase)
    }

    // Duplicate commit from the same C1 is a no-op and should not error.
    if err := st.Process(Message{ID: "blk1", From: "C1", Type: MsgCommit, Height: 7, Round: 1}); err != nil {
        t.Fatalf("duplicate commit must not error: %v", err)
    }
    if st.Phase != "commit" {
        t.Fatalf("duplicate commit must not change phase, got %q", st.Phase)
    }
}

