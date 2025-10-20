package qbft

import "testing"

// Prepare votes from different sources (From) should count once each and
// reach prepared at the threshold; any duplicates or additional prepares
// must not trigger another advance.
func TestState_Prepare_MultiSource_Dedup_Once(t *testing.T) {
    st := &State{Leader: "L"}

    // Establish proposal id and enter preprepared.
    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 11, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("want preprepared, got %q", st.Phase) }

    // Different From: P1 then P2 -> reach prepared.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 11, Round: 1}); err != nil {
        t.Fatalf("prepare P1: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("still preprepared after first prepare, got %q", st.Phase) }

    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 11, Round: 1}); err != nil {
        t.Fatalf("prepare P2: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("want prepared after second distinct prepare, got %q", st.Phase) }

    // Duplicate or extra prepares must not change phase further.
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 11, Round: 1}); err != nil {
        t.Fatalf("duplicate prepare must not error: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("phase changed unexpectedly after duplicate, got %q", st.Phase) }

    if err := st.Process(Message{ID: "blk1", From: "P3", Type: MsgPrepare, Height: 11, Round: 1}); err != nil {
        t.Fatalf("extra prepare must not error: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("phase changed unexpectedly after extra prepare, got %q", st.Phase) }
}

// Out-of-order prepares from different sources should still reach prepared
// once, and duplicates must not trigger another advance.
func TestState_Prepare_MultiSource_OutOfOrder(t *testing.T) {
    st := &State{Leader: "L"}

    if err := st.Process(Message{ID: "blk1", From: "L", Type: MsgPreprepare, Height: 12, Round: 0}); err != nil {
        t.Fatalf("preprepare: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("want preprepared, got %q", st.Phase) }

    // P2 arrives before P1.
    if err := st.Process(Message{ID: "blk1", From: "P2", Type: MsgPrepare, Height: 12, Round: 1}); err != nil {
        t.Fatalf("prepare P2: %v", err)
    }
    if st.Phase != "preprepared" { t.Fatalf("still preprepared after first prepare, got %q", st.Phase) }

    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 12, Round: 1}); err != nil {
        t.Fatalf("prepare P1: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("want prepared after second distinct prepare, got %q", st.Phase) }

    // Duplicate should not change phase further.
    if err := st.Process(Message{ID: "blk1", From: "P1", Type: MsgPrepare, Height: 12, Round: 1}); err != nil {
        t.Fatalf("duplicate prepare must not error: %v", err)
    }
    if st.Phase != "prepared" { t.Fatalf("phase changed unexpectedly after duplicate, got %q", st.Phase) }
}

