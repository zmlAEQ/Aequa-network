package p2p

import (
    "context"
    "fmt"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

type Service struct{
    mgr   *Manager
    gate  Gate
    rman  *ResourceManager
    hook  Hook
}

func New() *Service { return &Service{ mgr: NewManager(), gate: AllowAllGate{}, rman: NewResourceManager(DefaultResourceLimits()), hook: NopHook{} } }
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

// NewWithOpts allows tests to inject gate/resource/hook.
func NewWithOpts(mgr *Manager, gate Gate, rman *ResourceManager, hook Hook) *Service {
    if mgr == nil { mgr = NewManager() }
    if gate == nil { gate = AllowAllGate{} }
    if rman == nil { rman = NewResourceManager(DefaultResourceLimits()) }
    if hook == nil { hook = NopHook{} }
    return &Service{mgr: mgr, gate: gate, rman: rman, hook: hook}
}

// Connect tries to admit and register a peer according to gate and resources.
func (s *Service) Connect(id PeerID) error {
    labels := map[string]string{"result":"allowed"}
    if !s.gate.Allow(id) {
        labels["result"] = "denied"
        metrics.Inc("p2p_conn_attempts_total", labels)
        return fmt.Errorf("peer denied")
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