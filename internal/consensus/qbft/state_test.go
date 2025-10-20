package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestState_Process_IncrementsCounterAndUpdatesCoords(t *testing.T) {
    metrics.Reset()
    st := &State{}
    msg := Message{ID: "1", From: "p", Type: MsgPrepare, Height: 7, Round: 2}
    if err := st.Process(msg); err != nil {
        t.Fatalf("process: %v", err)
    }
    if st.Height != 7 || st.Round != 2 || st.Phase != "prepare" {
        t.Fatalf("state mismatch: %+v", *st)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_state_transitions_total{type="prepare"} 1`) {
        t.Fatalf("missing transitions counter: %q", dump)
    }
}

func TestState_Process_PhaseMapping(t *testing.T) {
    st := &State{}
    _ = st.Process(Message{Type: MsgPreprepare})
    if st.Phase != "preprepare" { t.Fatalf("phase: %s", st.Phase) }
    _ = st.Process(Message{Type: MsgCommit})
    if st.Phase != "commit" { t.Fatalf("phase: %s", st.Phase) }
}