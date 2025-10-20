package consensus

import (
    "context"
    "sync/atomic"
    "testing"
    "time"

    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
    "github.com/zmlAEQ/Aequa-network/pkg/bus"
)

type stubVerifier struct{ calls int32 }

func (s *stubVerifier) Verify(msg qbft.Message) error {
    atomic.AddInt32(&s.calls, 1)
    return nil
}

// Ensure that when an event is consumed, the verifier is invoked.
func TestService_InvokesVerifierOnEvent(t *testing.T) {
    b := bus.New(4)
    s := NewWithSub(b.Subscribe())
    sv := &stubVerifier{}
    s.SetVerifier(sv)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := s.Start(ctx); err != nil { t.Fatalf("start: %v", err) }

    b.Publish(ctx, bus.Event{Kind: bus.KindDuty})

    // Wait briefly to allow processing
    time.Sleep(30 * time.Millisecond)

    if atomic.LoadInt32(&sv.calls) == 0 {
        t.Fatalf("expected verifier to be called")
    }
}

