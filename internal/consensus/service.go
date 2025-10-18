package consensus

import (
    "context"
    "modular-dvt-engine/pkg/lifecycle"
    "modular-dvt-engine/pkg/logger"
)

type Service struct{}

func New() *Service { return &Service{} }
func (s *Service) Name() string { return "consensus" }
func (s *Service) Start(ctx context.Context) error { logger.Info("consensus start (stub)\n"); return nil }
func (s *Service) Stop(ctx context.Context) error  { logger.Info("consensus stop (stub)\n"); return nil }

var _ lifecycle.Service = (*Service)(nil)
