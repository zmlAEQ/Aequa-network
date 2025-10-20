package qbft

import "testing"

// Ensure that a valid preprepare from the expected leader moves phase to preprepared.
func TestState_Preprepare_LeaderOK(t *testing.T) {
    st := &State{Leader: "leader1"}
    msg := Message{ID: "pp", From: "leader1", Type: MsgPreprepare, Height: 10, Round: 0}
    if err := st.Process(msg); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if st.Phase != "preprepared" {
        t.Fatalf("want preprepared, got %q", st.Phase)
    }
}