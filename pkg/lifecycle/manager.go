package lifecycle

import (
    "context"
    "sync"
)

type Service interface {
    Name() string
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type Manager struct {
    svcs []Service
}

func New() *Manager { return &Manager{} }
func (m *Manager) Add(s Service) { m.svcs = append(m.svcs, s) }

func (m *Manager) StartAll(ctx context.Context) error {
    for _, s := range m.svcs {
        if err := s.Start(ctx); err != nil { return err }
    }
    return nil
}

func (m *Manager) StopAll(ctx context.Context) error {
    var wg sync.WaitGroup
    for _, s := range m.svcs {
        svc := s
        wg.Add(1)
        go func() { _ = svc.Stop(ctx); wg.Done() }()
    }
    wg.Wait()
    return nil
}
