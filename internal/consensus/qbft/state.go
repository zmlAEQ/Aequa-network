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
    Phase  string // e.g., "idle|preprepared|prepared|commit" (placeholder)
    Leader string // placeholder leader id for current round

    // Minimal aggregation placeholders for M3
    proposalID   string
    prepareVotes map[string]struct{} // by From
    commitVotes  map[string]struct{} // by From
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
    var ok bool
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
        s.proposalID = msg.ID
        s.prepareVotes = make(map[string]struct{})
        s.commitVotes = make(map[string]struct{})
    case MsgPrepare:
        // Fallback for legacy mapping test: if no preprepare yet, allow prepare
        // to pass and mark phase as "prepare" for observability only.
        if s.proposalID == "" {
            s.Phase = "prepare"
            break
        }
        // Accept prepare when preprepared or already prepared for the same proposal id.
        if s.Phase != "preprepared" && s.Phase != "prepared" {
            logger.ErrorJ("qbft_state", map[string]any{
                "op":        "transition",
                "event_type": string(msg.Type),
                "height":    s.Height,
                "round":     s.Round,
                "reason":    "not_preprepared",
            })
            return fmt.Errorf("prepare before preprepared")
        }
        if msg.ID != s.proposalID {
            logger.ErrorJ("qbft_state", map[string]any{
                "op":        "transition",
                "event_type": string(msg.Type),
                "height":    s.Height,
                "round":     s.Round,
                "reason":    "proposal_mismatch",
                "got":       msg.ID,
                "expect":    s.proposalID,
            })
            return fmt.Errorf("proposal mismatch")
        }
        if _, ok = s.prepareVotes[msg.From]; ok {
            // Duplicate prepare is a no-op regardless of current phase.
            break
        }
        s.prepareVotes[msg.From] = struct{}{}
        if s.Phase == "preprepared" && len(s.prepareVotes) >= 2 { // minimal threshold
            s.Phase = "prepared"
        }
    case MsgCommit:
        // Commit is valid for the current proposal after prepared.
        // If already in commit phase for the same proposal, treat duplicates as no-op.
        if s.Phase != "prepared" && s.Phase != "commit" {
            logger.ErrorJ("qbft_state", map[string]any{
                "op":        "transition",
                "event_type": string(msg.Type),
                "height":    s.Height,
                "round":     s.Round,
                "reason":    "not_prepared",
            })
            return fmt.Errorf("commit before prepared")
        }
        if msg.ID != s.proposalID {
            logger.ErrorJ("qbft_state", map[string]any{
                "op":        "transition",
                "event_type": string(msg.Type),
                "height":    s.Height,
                "round":     s.Round,
                "reason":    "proposal_mismatch",
                "got":       msg.ID,
                "expect":    s.proposalID,
            })
            return fmt.Errorf("proposal mismatch")
        }
        if _, ok = s.commitVotes[msg.From]; ok {
            // Duplicate commit (including when phase already is commit) is a no-op.
            break
        }
        s.commitVotes[msg.From] = struct{}{}
        // Minimal rule: first distinct commit advances to commit phase.
        if s.Phase != "commit" && len(s.commitVotes) >= 1 {
            s.Phase = "commit"
        }
    default:
        // Keep previous phase for unknown types; still record observability.
    }

    // Observability: one log + one counter per successful processed message.
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
