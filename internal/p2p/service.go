package p2p

import (
    "context"
    "time"
    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

type Service struct{}

func New() *Service { return &Service{} }
func (s *Service) Name() string { return "p2p" }
func (s *Service) Start(ctx context.Context) error {
    begin := time.Now()
    dur := time.Since(begin).Milliseconds()
    logger.InfoJ("service_op", map[string]any{"service":"p2p", "op":"start", "result":"ok", "latency_ms": dur})
    metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"start"}, float64(dur))
    return nil
}
func (s *Service) Stop(ctx context.Context) error  {
    begin := time.Now()
    dur := time.Since(begin).Milliseconds()
    logger.InfoJ("service_op", map[string]any{"service":"p2p", "op":"stop", "result":"ok", "latency_ms": dur})
    metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"stop"}, float64(dur))
    return nil
}

var _ lifecycle.Service = (*Service)(nil)


