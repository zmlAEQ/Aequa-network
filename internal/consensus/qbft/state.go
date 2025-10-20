package qbft

import ("fmt"`r`n    "github.com/zmlAEQ/Aequa-network/pkg/logger"`r`n    "github.com/zmlAEQ/Aequa-network/pkg/metrics"`r`n)

// State represents a minimal QBFT state snapshot.
// This is a skeleton for M3: it carries only coordinates and a textual phase.
type State struct {`r`n    Height uint64`r`n    Round  uint64`r`n    Phase  string // e.g., "idle|preprepared|prepare|commit" (placeholder)`r`n    Leader string // placeholder leader id for current round`r`n}

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
    case MsgPreprepare:`r`n        if s.Leader != "" && msg.From != s.Leader {`r`n            logger.ErrorJ("qbft_state", map[string]any{`r`n                "op": "transition",`r`n                "event_type": string(msg.Type),`r`n                "height": s.Height,`r`n                "round": s.Round,`r`n                "reason": "unauthorized_leader",`r`n                "from": msg.From,`r`n                "expect": s.Leader,`r`n            })`r`n            return fmt.Errorf("unauthorized leader")`r`n        }`r`n        s.Phase = "preprepared"`r`n    case MsgPrepare:
        s.Phase = "prepare"
    case MsgCommit:
        s.Phase = "commit"
    default:
        // Keep previous phase for unknown types; still record observability.
    }

    // Observability: one log + one counter per processed message.
    logger.InfoJ("qbft_state", map[string]any{
        "op": "transition",
        "event_type": string(msg.Type),
        "height": s.Height,
        "round": s.Round,
        "phase": s.Phase,
    })
    metrics.Inc("qbft_state_transitions_total", map[string]string{"type": string(msg.Type)})
    return nil
}

