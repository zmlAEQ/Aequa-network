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

type BasicVerifier struct{ replay *AntiReplay; minHeight uint64 }

func NewBasicVerifier() *BasicVerifier { return &BasicVerifier{replay: NewAntiReplay()} }\n\nfunc (v *BasicVerifier) SetMinHeight(h uint64) { v.minHeight = h }\nfunc (v *BasicVerifier) SetRoundWindow(w uint64) { v.roundWindow = w }\nfunc (v *BasicVerifier) SetAllowed(ids ...string) { if v.allowed == nil { v.allowed = map[string]struct{}{} }; for _, id := range ids { v.allowed[id] = struct{}{} } }\n\nfunc validType(t Type) bool {\n    switch t {\n    case MsgPreprepare, MsgPrepare, MsgCommit:\n        return true\n    default:\n        return false\n    }\n}\n\nfunc (v *BasicVerifier) SetMinHeight(h uint64) { v.minHeight = h }\n\nfunc validType(t Type) bool {\n    switch t {\n    case MsgPreprepare, MsgPrepare, MsgCommit:\n        return true\n    default:\n        return false\n    }\n}

func (v *BasicVerifier) Verify(msg Message) error {
    labels := map[string]string{"type": string(msg.Type)}
    // basic structural checks\n    if msg.ID == "" || msg.From == "" || msg.Type == "" || !validType(msg.Type) {\n        metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"error"})\n        logger.ErrorJ("qbft_verify", map[string]any{"result":"error", "reason":"invalid", "type": string(msg.Type)})\n        return fmt.Errorf("invalid message")\n    }\n    // sender whitelist (optional)\n    if len(v.allowed) > 0 { if _, ok := v.allowed[msg.From]; !ok { metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"unauthorized"}); logger.ErrorJ("qbft_verify", map[string]any{"result":"unauthorized", "from": msg.From, "type": string(msg.Type)}); return fmt.Errorf("unauthorized") } }\n    // signature shape placeholder\n    if msg.Sig != nil && len(msg.Sig) > 0 && len(msg.Sig) < 32 { metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"sig_invalid"}); logger.ErrorJ("qbft_verify", map[string]any{"result":"sig_invalid", "type": string(msg.Type)}); return fmt.Errorf("sig invalid") }\n    // height window (old height)\n    if v.minHeight > 0 && msg.Height < v.minHeight { metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"old"}); logger.ErrorJ("qbft_verify", map[string]any{"result":"old", "height": msg.Height, "min": v.minHeight, "type": string(msg.Type)}); return fmt.Errorf("old height") }\n    // round window (optional upper bound)\n    if v.roundWindow > 0 && msg.Round > v.roundWindow { metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"round_oob"}); logger.ErrorJ("qbft_verify", map[string]any{"result":"round_oob", "round": msg.Round, "max": v.roundWindow, "type": string(msg.Type)}); return fmt.Errorf("round out of bound") }
        // height window (old height)\n    if v.minHeight > 0 && msg.Height < v.minHeight {\n        metrics.Inc("qbft_msg_verified_total", map[string]string{"result":"old"})\n        logger.ErrorJ("qbft_verify", map[string]any{"result":"old", "height": msg.Height, "min": v.minHeight, "type": string(msg.Type)})\n        return fmt.Errorf("old height")\n    }\n    // anti-replay
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