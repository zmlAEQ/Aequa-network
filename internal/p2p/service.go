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

type Service struct{`r`n    mgr   *Manager`r`n    gate  Gate`r`n    rman  *ResourceManager`r`n    hook  Hook`r`n    dkgv  dkg.Verifier`r`n    cfg   Config`r`n}

func New() *Service { return &Service{ mgr: NewManager(), gate: AllowAllGate{}, rman: NewResourceManager(DefaultResourceLimits()), hook: LogHook{}, dkgv: dkg.NoopVerifier{}, cfg: DefaultConfig() } }, rman: NewResourceManager(DefaultResourceLimits()), hook: LogHook{}, dkgv: dkg.NoopVerifier{} } }
func (s *Service) Name() string { return "p2p" }
func (s *Service) Start(ctx context.Context) error {`r`n    begin := time.Now()`r`n    // Config validation (fail-fast)`r`n    if err := s.cfg.Validate(s.dkgv != nil); err != nil {`r`n        metrics.Inc("p2p_config_checks_total", map[string]string{"result":"error"})`r`n        logger.ErrorJ("p2p_config", map[string]any{"result":"error", "err": err.Error()})`r`n        dur := time.Since(begin).Milliseconds()`r`n        logger.ErrorJ("service_op", map[string]any{"service":"p2p", "op":"start", "result":"error", "latency_ms": dur})`r`n        metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"start"}, float64(dur))`r`n        return err`r`n    }`r`n    metrics.Inc("p2p_config_checks_total", map[string]string{"result":"ok"})`r`n    // DKG/cluster-lock verification (fail-fast)`r`n    if s.dkgv != nil {`r`n        if err := s.dkgv.VerifyCluster(); err != nil {`r`n            logger.ErrorJ("p2p_dkg_cluster", map[string]any{"result":"error", "err": err.Error()})`r`n            metrics.Inc("p2p_dkg_cluster_checks_total", map[string]string{"result":"error"})`r`n            dur := time.Since(begin).Milliseconds()`r`n            logger.ErrorJ("service_op", map[string]any{"service":"p2p", "op":"start", "result":"error", "latency_ms": dur})`r`n            metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"start"}, float64(dur))`r`n            return err`r`n        } else {`r`n            logger.InfoJ("p2p_dkg_cluster", map[string]any{"result":"ok"})`r`n            metrics.Inc("p2p_dkg_cluster_checks_total", map[string]string{"result":"ok"})`r`n        }`r`n    }`r`n    dur := time.Since(begin).Milliseconds()`r`n    logger.InfoJ("service_op", map[string]any{"service":"p2p", "op":"start", "result":"ok", "latency_ms": dur})`r`n    metrics.ObserveSummary("service_op_ms", map[string]string{"service":"p2p", "op":"start"}, float64(dur))`r`n    return nil`r`n}
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
    return &Service{mgr: mgr, gate: gate, rman: rman, hook: hook, dkgv: dkg.NoopVerifier{}}
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
// SetConfig injects a validated P2P config.
func (s *Service) SetConfig(c Config) { s.cfg = c }
