package qbft

import "testing"

// Duplicate Prepare from the same From for the same proposal ID must not
// increase the effective vote count nor advance the phase. Processing the
// duplicate should be a no-op and return no error.
func TestState_Prepare_Dedup_SameFrom_DoesNotAdvanceOrCount(t *testing.T) {
    st := &State{Leader: "L"}

    // Move to preprepared with proposal id "blk1".
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 5, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("want preprepared, got %q", st.Phase)
    }

    // First prepare from P1 should not yet advance to prepared.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 5, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("still expect preprepared after first prepare, got %q", st.Phase)
    }

    // Duplicate prepare from the same P1 (same proposal id) is ignored.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 5, Round: 1}); err != nil {
        t.Fatalf("duplicate prepare must not error: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("duplicate prepare must not advance phase, got %q", st.Phase)
    }

    // A different validator P2 should now reach the threshold and advance.
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 5, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" {
        t.Fatalf("want prepared after distinct second prepare, got %q", st.Phase)
    }
}

