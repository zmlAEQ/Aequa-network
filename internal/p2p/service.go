package p2p

import (
    "context"
    "fmt"
    "time"

    "github.com/zmlAEQ/Aequa-network/internal/dkg"
    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

type Service struct{
    mgr   *Manager
    gate  Gate
    rman  *ResourceManager
    hook  Hook
    dkgv  dkg.Verifier
    cfg   Config
}

func New() *Service { return &Service{ mgr: NewManager(), gate: AllowAllGate{}, rman: NewResourceManager(DefaultResourceLimits()), hook: LogHook{}, dkgv: dkg.NoopVerifier{}, cfg: DefaultConfig() } }
func (s *Service) Name() string { return "p2p" }

// SetConfig injects a validated P2P config.
func (s *Service) SetConfig(c Config) { s.cfg = c }

func (s *Service) Start(ctx context.Context) error {
    begin := time.Now()
    // Config validation (fail-fast)
    if err := s.cfg.Validate(s.dkgv != nil); err != nil {
        metrics.Inc("p2p_config_checks_total", map[string]string{"result":"error"})
        logger.ErrorJ("p2p_config", map[string]any{"result":"error", "err": err.Error()})
        dur := time.Since(begin).Milliseconds()
        logger.ErrorJ("service_op", map[string]any{"service":"p2p", "op":"start", "result":"error", "latency_ms": dur})
        metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"start"}, float64(dur))
        return err
    }
    metrics.Inc("p2p_config_checks_total", map[string]string{"result":"ok"})

    // DKG/cluster-lock verification (fail-fast)
    if s.dkgv != nil {
        if err := s.dkgv.VerifyCluster(); err != nil {
            logger.ErrorJ("p2p_dkg_cluster", map[string]any{"result":"error", "err": err.Error()})
            metrics.Inc("p2p_dkg_cluster_checks_total", map[string]string{"result":"error"})
            dur := time.Since(begin).Milliseconds()
            logger.ErrorJ("service_op", map[string]any{"service":"p2p", "op":"start", "result":"error", "latency_ms": dur})
            metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"start"}, float64(dur))
            return err
        } else {
            logger.InfoJ("p2p_dkg_cluster", map[string]any{"result":"ok"})
            metrics.Inc("p2p_dkg_cluster_checks_total", map[string]string{"result":"ok"})
        }
    }

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

// NewWithOpts allows tests to inject gate/resource/hook.
func NewWithOpts(mgr *Manager, gate Gate, rman *ResourceManager, hook Hook) *Service {
    if mgr == nil { mgr = NewManager() }
    if gate == nil { gate = AllowAllGate{} }
    if rman == nil { rman = NewResourceManager(DefaultResourceLimits()) }
    if hook == nil { hook = NopHook{} }
    return &Service{mgr: mgr, gate: gate, rman: rman, hook: hook, dkgv: dkg.NoopVerifier{}, cfg: DefaultConfig()}
}

// SetDKG allows tests or wiring to inject a DKG/cluster-lock verifier.
func (s *Service) SetDKG(v dkg.Verifier) { s.dkgv = v }

// Connect tries to admit and register a peer according to gate and resources.
func (s *Service) Connect(id PeerID) error {
    labels := map[string]string{"result":"allowed"}
    if rg, ok := any(s.gate).(ReasonedGate); ok {
        if ok2, reason := rg.AllowWithReason(id); !ok2 {
            if reason == "" { reason = "denied" }
            labels["result"] = reason
            metrics.Inc("p2p_conn_attempts_total", labels)
            return fmt.Errorf("peer denied")
        }
    } else {
        if !s.gate.Allow(id) {
            labels["result"] = "denied"
            metrics.Inc("p2p_conn_attempts_total", labels)
            return fmt.Errorf("peer denied")
        }
    }
    if s.dkgv != nil && !s.dkgv.AllowPeer(string(id)) {
        labels["result"] = "dkg_denied"
        metrics.Inc("p2p_conn_attempts_total", labels)
        logger.ErrorJ("p2p_dkg_gate", map[string]any{"peer_id": string(id), "result":"denied"})
        return fmt.Errorf("dkg denied")
    }
    if !s.rman.TryOpen() {
        labels["result"] = "limited"
        metrics.Inc("p2p_conn_attempts_total", labels)
        return fmt.Errorf("resource limit")
    }
    s.mgr.AddPeer(id)
    metrics.Inc("p2p_conn_attempts_total", labels)
    s.hook.OnPeerJoin(string(id))
    return nil
}

// Disconnect unregisters a peer and releases resources.
func (s *Service) Disconnect(id PeerID) {
    s.mgr.RemovePeer(id)
    s.rman.Close()
    s.hook.OnPeerLeave(string(id))
}