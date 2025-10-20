package qbft

import "testing"

// In a simplified model, after reaching a prepare threshold, phase moves to prepared.
func TestState_Prepared_AfterThreshold(t *testing.T) {
    st := &State{Leader: "L", PrepareThreshold: 2}
    // Accept preprepare from leader and establish proposal id
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 5, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("want preprepared, got %q", st.Phase) }

    // First prepare from P1 (same proposal id)
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 5, Round: 1}); err != nil {
        t.Fatalf("prepare1: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("should still be preprepared after first prepare, got %q", st.Phase) }

    // Second prepare from P2 reaches threshold -> prepared
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 5, Round: 1}); err != nil {
        t.Fatalf("prepare2: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("want prepared, got %q", st.Phase) }
}