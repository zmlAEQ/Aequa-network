package consensus

import (
    "context"
    "testing"
    "time"

    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
    "github.com/zmlAEQ/Aequa-network/pkg/bus"
)

type okVerifier struct{}
func (okVerifier) Verify(msg qbft.Message) error { return nil }

// Ensure that after a verified message, the qbft.State processor updates Phase as expected.
func TestService_StateProcessor_PhaseUpdatesOnVerifiedMessage(t *testing.T) {
    b := bus.New(8)
    s := NewWithSub(b.Subscribe())

    st := &qbft.State{}
    s.SetProcessor(st)
    s.SetVerifier(okVerifier{})

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := s.Start(ctx); err != nil { t.Fatalf("start: %v", err) }

    // Publish a legal event; adapter maps to MsgPrepare with given Height/Round
    b.Publish(ctx, bus.Event{Kind: bus.KindDuty, Height: 11, Round: 1, TraceID: "tid"})

    time.Sleep(50 * time.Millisecond)

    if got := st.Phase; got != "prepare" {
        t.Fatalf("phase not updated: got %q, want %q", got, "prepare")
    }
    if st.Height != 11 || st.Round != 1 {
        t.Fatalf("coords not updated: got (%d,%d)", st.Height, st.Round)
    }
}
