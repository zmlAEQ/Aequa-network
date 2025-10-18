package lifecycle

import (
    "context"
    "errors"
)

type Service interface {
    Name() string
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type Manager struct { svcs []Service }

func New() *Manager { return &Manager{} }
func (m *Manager) Add(s Service) { m.svcs = append(m.svcs, s) }

// StartAll starts services in order. On failure it stops the already started
// services in reverse order and returns a joined error containing the original
// start error and any rollback errors.
func (m *Manager) StartAll(ctx context.Context) error {
    started := 0
    for i, s := range m.svcs {
        if err := s.Start(ctx); err != nil {
            // rollback previously started services in reverse
            var merr error = err
            for j := i - 1; j >= 0; j-- {
                if rerr := m.svcs[j].Stop(ctx); rerr != nil { merr = errors.Join(merr, rerr) }
            }
            return merr
        }
        started++
    }
    _ = started
    return nil
}

// StopAll stops services in reverse order. It attempts to stop all and returns
// a joined error if any stop fails.
func (m *Manager) StopAll(ctx context.Context) error {
    var merr error
    for i := len(m.svcs) - 1; i >= 0; i-- {
        if err := m.svcs[i].Stop(ctx); err != nil { merr = errors.Join(merr, err) }
    }
    return merr
}
