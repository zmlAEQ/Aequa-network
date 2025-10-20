package qbft

import (
    "fmt"
    "sync"

    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

type Verifier interface {
    Verify(msg Message) error
}

type AntiReplay struct {
    mu   sync.Mutex
    seen map[string]struct{}
}

func NewAntiReplay() *AntiReplay { return &AntiReplay{seen: make(map[string]struct{})} }

// Seen returns true if id already seen; otherwise records and returns false.
func (r *AntiReplay) Seen(id string) bool {
    if id == "" { return false }
    r.mu.Lock(); defer r.mu.Unlock()
    if _, ok := r.seen[id]; ok { return true }
    r.seen[id] = struct{}{}
    return false
}

type BasicVerifier struct{ replay *AntiReplay }

func NewBasicVerifier() *BasicVerifier { return &BasicVerifier{replay: NewAntiReplay()} }

func (v *BasicVerifier) Verify(msg Message) error {
    labels := map[string]string{"type": string(msg.Type)}
    // basic structural checks
    if msg.ID == "" || msg.From == "" || msg.Type == "" {
        metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"error"})
        logger.ErrorJ("qbft_verify", map[string]any{"result":"error", "reason":"invalid", "type": string(msg.Type)})
        return fmt.Errorf("invalid message")
    }
    // anti-replay
    if v.replay != nil && v.replay.Seen(msg.ID) {
        metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"replay"})
        logger.ErrorJ("qbft_verify", map[string]any{"result":"replay", "id": msg.ID, "type": string(msg.Type)})
        return fmt.Errorf("replay")
    }
    // ok
    metrics.Inc("qbft_msg_verified_total", labels)
    logger.InfoJ("qbft_verify", map[string]any{"result":"ok", "id": msg.ID, "type": string(msg.Type)})
    return nil
}