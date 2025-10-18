package consensus

import (
    "context"
    "github.com/zmlAEQ/Aequa-network/pkg/bus"
    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
)

type Service struct{ sub bus.Subscriber }

func New() *Service { return &Service{} }
func NewWithSub(sub bus.Subscriber) *Service { return &Service{sub: sub} }
func (s *Service) Name() string { return "consensus" }

func (s *Service) Start(ctx context.Context) error {
    if s.sub == nil {
        logger.Info("consensus start (stub)")
        return nil
    }
    go func() {
        for {
            select {
            case ev := <-s.sub:
                logger.InfoJ("consensus_recv", map[string]any{"kind": string(ev.Kind)})
            case <-ctx.Done():
                return
            }
        }
    }()
    return nil
}

func (s *Service) Stop(ctx context.Context) error  { logger.Info("consensus stop (stub)"); return nil }

var _ lifecycle.Service = (*Service)(nil)

