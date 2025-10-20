package consensus

import (
    "context"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/bus"
    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
    "github.com/zmlAEQ/Aequa-network/internal/state"
)

type Service struct{ sub bus.Subscriber; v qbft.Verifier; store state.Store }

func New() *Service { return &Service{} }
func NewWithSub(sub bus.Subscriber) *Service { return &Service{sub: sub} }
func (s *Service) Name() string { return "consensus" }

// SetVerifier allows tests/wiring to inject a qbft.Verifier. If nil, a BasicVerifier is instantiated on start.
func (s *Service) SetVerifier(v qbft.Verifier) { s.v = v }

// SetStore allows tests/wiring to inject a StateDB store. If nil, a MemoryStore is instantiated on start.
func (s *Service) SetStore(st state.Store) { s.store = st }

func (s *Service) Start(ctx context.Context) error {
    if s.sub == nil {
        logger.Info("consensus start (stub)")
        return nil
    }
    if s.v == nil { s.v = qbft.NewBasicVerifier() }
    if s.store == nil { s.store = state.NewMemoryStore() }
    if ls, err := s.store.LoadLastState(ctx); err != nil {
        logger.InfoJ("consensus_state", map[string]any{"op":"load", "result":"miss", "err": err.Error()})
    } else {
        logger.InfoJ("consensus_state", map[string]any{"op":"load", "result":"ok", "height": ls.Height, "round": ls.Round})
    }
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
                // Map event to qbft message via adapter (stub mapping only)
                msg := MapEventToQBFT(ev)
                _ = s.v.Verify(msg)
                // Persist last state (stub): log only; ignore control flow
                if err := s.store.SaveLastState(ctx, state.LastState{Height: msg.Height, Round: msg.Round}); err != nil {
                    logger.ErrorJ("consensus_state", map[string]any{"op":"save", "result":"error", "err": err.Error()})
                } else {
                    logger.InfoJ("consensus_state", map[string]any{"op":"save", "result":"ok", "height": msg.Height, "round": msg.Round})
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    return nil
}

func (s *Service) Stop(ctx context.Context) error  { logger.Info("consensus stop (stub)"); return nil }

var _ lifecycle.Service = (*Service)(nil)


