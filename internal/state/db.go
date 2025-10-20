package state

import (
    "context"
    "errors"
    "sync"
)

// LastState represents the latest consensus coordinates persisted by the node.
// It is intentionally small and generic for M3 minimal persistence.
type LastState struct {
    Height uint64
    Round  uint64
}

// ErrNotFound is returned when a requested state is not available.
var ErrNotFound = errors.New("not found")

// Store defines the minimal interface for persisting and loading the last state.
// Implementations should be concurrency-safe.
type Store interface {
    SaveLastState(ctx context.Context, s LastState) error
    LoadLastState(ctx context.Context) (LastState, error)
    Close() error
}

// MemoryStore is a minimal in-memory implementation of Store.
// It is intended as a stub for wiring and tests in M3 and is not durable.
type MemoryStore struct {
    mu   sync.RWMutex
    have bool
    last LastState
}

// NewMemoryStore constructs a new empty MemoryStore.
func NewMemoryStore() *MemoryStore { return &MemoryStore{} }

// SaveLastState stores the provided state atomically.
func (m *MemoryStore) SaveLastState(_ context.Context, s LastState) error {
    m.mu.Lock()
    m.last = s
    m.have = true
    m.mu.Unlock()
    return nil
}

// LoadLastState returns the last stored state, or ErrNotFound if none.
func (m *MemoryStore) LoadLastState(_ context.Context) (LastState, error) {
    m.mu.RLock()
    have, s := m.have, m.last
    m.mu.RUnlock()
    if !have { return LastState{}, ErrNotFound }
    return s, nil
}

// Close implements Store. For MemoryStore it is a no-op.
func (m *MemoryStore) Close() error { return nil }


