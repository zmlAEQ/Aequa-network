package p2p

import (
    "sync"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// PeerID is a lightweight identifier for a peer (placeholder for libp2p peer.ID).
type PeerID string

// Message represents a broadcast payload.
type Message struct {
    From    PeerID
    Payload []byte
    TraceID string
}

// Manager provides minimal peer tracking and broadcast metrics without real networking.
type Manager struct {
    mu    sync.RWMutex
    peers map[PeerID]struct{}
}

func NewManager() *Manager { return &Manager{peers: make(map[PeerID]struct{})} }

func (m *Manager) AddPeer(id PeerID) {
    m.mu.Lock(); m.peers[id] = struct{}{}; m.mu.Unlock()
    metrics.Inc("p2p_peers_added_total", nil)
}

func (m *Manager) RemovePeer(id PeerID) {
    m.mu.Lock(); delete(m.peers, id); m.mu.Unlock()
    metrics.Inc("p2p_peers_removed_total", nil)
}

// Broadcast simulates broadcasting a message to current peers and records metrics.
func (m *Manager) Broadcast(msg Message) int {
    begin := time.Now()
    m.mu.RLock(); n := len(m.peers); m.mu.RUnlock()
    metrics.Inc("p2p_messages_total", map[string]string{"kind":"broadcast"})
    metrics.ObserveSummary("p2p_broadcast_ms", map[string]string{"kind":"broadcast"}, float64(time.Since(begin).Milliseconds()))
    logger.InfoJ("p2p_broadcast", map[string]any{"peers": n, "latency_ms": time.Since(begin).Milliseconds(), "trace_id": msg.TraceID})
    return n
}