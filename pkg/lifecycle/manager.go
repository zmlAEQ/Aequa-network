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

type Manager struct {
    svcs []Service
}

func New() *Manager { return &Manager{} }
func (m *Manager) Add(s Service) { m.svcs = append(m.svcs, s) }

// StartAll 按顺序启动；若某服务启动失败，则按逆序回滚已启动服务，并聚合错误返回。
func (m *Manager) StartAll(ctx context.Context) error {
    var started []Service
    for _, s := range m.svcs {
        if err := s.Start(ctx); err != nil {
            agg := err
            for i := len(started) - 1; i >= 0; i-- {
                if stopErr := started[i].Stop(ctx); stopErr != nil {
                    agg = errors.Join(agg, stopErr)
                }
            }
            return agg
        }
        started = append(started, s)
    }
    return nil
}

// StopAll 按逆序关闭所有服务，并聚合所有停止错误。
func (m *Manager) StopAll(ctx context.Context) error {
    var agg error
    for i := len(m.svcs) - 1; i >= 0; i-- {
        if err := m.svcs[i].Stop(ctx); err != nil {
            agg = errors.Join(agg, err)
        }
    }
    return agg
}