package p2p

import (
    "context"
    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
)

type Service struct{}

func New() *Service { return &Service{} }
func (s *Service) Name() string { return "p2p" }
func (s *Service) Start(ctx context.Context) error { logger.Info("p2p start (stub)\n"); return nil }
func (s *Service) Stop(ctx context.Context) error  { logger.Info("p2p stop (stub)\n"); return nil }

var _ lifecycle.Service = (*Service)(nil)


