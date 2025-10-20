package state

import (
    "context"
    "testing"
)

func TestMemoryStore_NotFound(t *testing.T) {
    m := NewMemoryStore()
    if _, err := m.LoadLastState(context.Background()); err != ErrNotFound {
        t.Fatalf("expected ErrNotFound, got %v", err)
    }
}

func TestMemoryStore_SaveThenLoad(t *testing.T) {
    m := NewMemoryStore()
    want := LastState{Height: 123, Round: 4}
    if err := m.SaveLastState(context.Background(), want); err != nil {
        t.Fatalf("save: %v", err)
    }
    got, err := m.LoadLastState(context.Background())
    if err != nil {
        t.Fatalf("load: %v", err)
    }
    if got != want {
        t.Fatalf("mismatch: got %+v, want %+v", got, want)
    }
}