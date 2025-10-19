package p2p

import (
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
)

type Hook interface {
    OnPeerJoin(id string)
    OnPeerLeave(id string)
}

type NopHook struct{}
func (NopHook) OnPeerJoin(id string) {}
func (NopHook) OnPeerLeave(id string) {}

// LogHook emits audit logs and increments counters for peer join/leave.
type LogHook struct{}

func (LogHook) OnPeerJoin(id string) {
    metrics.Inc("p2p_peers_joined_total", nil)
    logger.InfoJ("p2p_peer", map[string]any{
        "event": "join",
        "peer_id": id,
        "latency_ms": 0,
        "result": "ok",
        "trace_id": "",
    })
}

func (LogHook) OnPeerLeave(id string) {
    metrics.Inc("p2p_peers_left_total", nil)
    logger.InfoJ("p2p_peer", map[string]any{
        "event": "leave",
        "peer_id": id,
        "latency_ms": 0,
        "result": "ok",
        "trace_id": "",
    })
}