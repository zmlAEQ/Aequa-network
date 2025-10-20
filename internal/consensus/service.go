package consensus

import (
    "context"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/bus"
    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
)

type Service struct{ sub bus.Subscriber; v qbft.Verifier }

func New() *Service { return &Service{} }
func NewWithSub(sub bus.Subscriber) *Service { return &Service{sub: sub} }
func (s *Service) Name() string { return "consensus" }

// SetVerifier allows tests/wiring to inject a qbft.Verifier. If nil, a BasicVerifier is instantiated on start.
func (s *Service) SetVerifier(v qbft.Verifier) { s.v = v }

func (s *Service) Start(ctx context.Context) error {
    if s.sub == nil {
        logger.Info("consensus start (stub)")
        return nil
    }
    if s.v == nil { s.v = qbft.NewBasicVerifier() }
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case ev := <-s.sub:
                start := time.Now()
                // Audit log + metrics for event intake
                dur := time.Since(start)
                logger.InfoJ("consensus_recv", map[string]any{"kind": string(ev.Kind), "trace_id": ev.TraceID, "result": "recv", "latency_ms": dur.Milliseconds()})
                metrics.Inc("consensus_events_total", map[string]string{"kind": string(ev.Kind)})
                metrics.ObserveSummary("consensus_proc_ms", map[string]string{"kind": string(ev.Kind)}, float64(dur.Milliseconds()))
                // Verify placeholder qbft message (stub). Do not alter control flow.
                msg := qbft.Message{
                    ID:      time.Now().Format("20060102T150405.000000000"),
                    From:    "consensus_stub",
                    Type:    qbft.MsgPrepare,
                    Height:  0,
                    Round:   1,
                    Payload: nil,
                    TraceID: ev.TraceID,
                    Sig:     nil,
                }
                _ = s.v.Verify(msg)
            case <-ctx.Done():
                return
            }
        }
    }()
    return nil
}

func (s *Service) Stop(ctx context.Context) error  { logger.Info("consensus stop (stub)"); return nil }

var _ lifecycle.Service = (*Service)(nil)


