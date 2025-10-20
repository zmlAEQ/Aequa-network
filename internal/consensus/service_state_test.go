package consensus

import (
    "context"
    "sync/atomic"
    "testing"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/bus"
    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
    "github.com/zmlAEQ/Aequa-network/internal/state"
)

type stubStore struct{
    loads int32
    saves int32
    last  state.LastState
}

func (s *stubStore) SaveLastState(_ context.Context, ls state.LastState) error {
    atomic.AddInt32(&s.saves, 1)
    s.last = ls
    return nil
}
func (s *stubStore) LoadLastState(_ context.Context) (state.LastState, error) {
    atomic.AddInt32(&s.loads, 1)
    return state.LastState{}, state.ErrNotFound
}
func (s *stubStore) Close() error { return nil }

type nopVerifier struct{}
func (nopVerifier) Verify(msg qbft.Message) error { return nil }

// Ensure consensus.Service loads at start and saves after processing an event.
func TestService_StateStore_LoadAndSaveCalled(t *testing.T) {
    b := bus.New(4)
    s := NewWithSub(b.Subscribe())
    st := &stubStore{}
    s.SetStore(st)
    s.SetVerifier(nopVerifier{}) // avoid affecting logs/metrics expectations

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := s.Start(ctx); err != nil { t.Fatalf("start: %v", err) }

    // One event triggers a save after verify
    b.Publish(ctx, bus.Event{Kind: bus.KindDuty})
    time.Sleep(30 * time.Millisecond)

    if atomic.LoadInt32(&st.loads) == 0 {
        t.Fatalf("expected LoadLastState to be called at start")
    }
    if atomic.LoadInt32(&st.saves) == 0 {
        t.Fatalf("expected SaveLastState to be called after event")
    }
}

