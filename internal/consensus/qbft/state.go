package qbft

import (
    "fmt"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// State represents a minimal QBFT state snapshot.
// This is a skeleton for M3: it carries only coordinates and a textual phase.
type State struct {
    Height uint64
    Round  uint64
    Phase  string // e.g., "idle|preprepared|prepare|commit" (placeholder)
    Leader string // placeholder leader id for current round
}

// Processor defines the minimal interface for driving state transitions.
type Processor interface {
    Process(msg Message) error
}

// Process triggers a placeholder state transition based on the incoming message.
// It does not enforce any real QBFT rules; it only updates coordinates,
// emits a log, and increments a Prometheus counter for observability.
func (s *State) Process(msg Message) error {
    // Lightweight, non-authoritative update of coordinates for visibility.
    s.Height = msg.Height
    s.Round = msg.Round
    switch msg.Type {
    case MsgPreprepare:
        // Placeholder leader validation: if Leader is set, only accept from that id
        if s.Leader != "" && msg.From != s.Leader {
            logger.ErrorJ("qbft_state", map[string]any{
                "op":        "transition",
                "event_type": string(msg.Type),
                "height":    s.Height,
                "round":     s.Round,
                "reason":    "unauthorized_leader",
                "from":      msg.From,
                "expect":    s.Leader,
            })
            return fmt.Errorf("unauthorized leader")
        }
        s.Phase = "preprepared"
    case MsgPrepare:
        s.Phase = "prepare"
    case MsgCommit:
        s.Phase = "commit"
    default:
        // Keep previous phase for unknown types; still record observability.
    }

    // Observability: one log + one counter per processed message.
    logger.InfoJ("qbft_state", map[string]any{
        "op":        "transition",
        "event_type": string(msg.Type),
        "height":    s.Height,
        "round":     s.Round,
        "phase":     s.Phase,
    })
    metrics.Inc("qbft_state_transitions_total", map[string]string{"type": string(msg.Type)})
    return nil
}