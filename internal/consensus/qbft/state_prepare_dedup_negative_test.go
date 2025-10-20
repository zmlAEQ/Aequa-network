package qbft

import "testing"

// Duplicate prepares from the same sender for the same proposal must not
// increase vote count nor advance phase to prepared. It may return nil error.
func TestState_Prepare_DuplicateFrom_DoesNotAdvance(t *testing.T) {
    st := &State{Leader: "L"}

    // Enter preprepared with proposal id blk1.
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 9, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("want preprepared, got %q", st.Phase) }

    // First prepare from P1 should not yet reach prepared (default thr=2).
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 9, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("after first prepare, phase=%q", st.Phase) }

    // Duplicate prepare from the same P1 must not advance phase.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 9, Round: 1}); err != nil {
        t.Fatalf("duplicate prepare unexpected error: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("duplicate must not advance: %q", st.Phase) }
}