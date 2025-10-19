package p2p

import (
    "sync/atomic"
)// Gate decides whether a peer is allowed to connect.
type Gate interface { Allow(id PeerID) bool }

// AllowAllGate admits all peers (default in minimal setup).
type AllowAllGate struct{}
func (AllowAllGate) Allow(id PeerID) bool { return true }

// StaticDenyGate denies all peers (used in tests).
type StaticDenyGate struct{}
func (StaticDenyGate) Allow(id PeerID) bool { return false }
// AllowListGate admits only peers present in its in-memory allowlist.
type AllowListGate struct{ allowed map[PeerID]struct{} }

// NewAllowListGate constructs an allowlist gate from the provided peer IDs.
func NewAllowListGate(ids ...PeerID) AllowListGate {
    m := make(map[PeerID]struct{}, len(ids))
    for _, id := range ids { m[id] = struct{}{} }
    return AllowListGate{allowed: m}
}

func (g AllowListGate) Allow(id PeerID) bool {
    _, ok := g.allowed[id]
    return ok
}
// ReasonedGate optionally returns a reason when denying a peer.
type ReasonedGate interface { AllowWithReason(id PeerID) (bool, string) }

// RateLimitGate allows only a fixed number of connections (process-wide stub).
type RateLimitGate struct{ remain int64 }

func NewRateLimitGate(limit int64) RateLimitGate { return RateLimitGate{remain: limit} }

func (g *RateLimitGate) Allow(id PeerID) bool {
    ok, _ := g.AllowWithReason(id)
    return ok
}

func (g *RateLimitGate) AllowWithReason(id PeerID) (bool, string) {
    for {
        r := atomic.LoadInt64(&g.remain)
        if r <= 0 { return false, "rate_limited" }
        if atomic.CompareAndSwapInt64(&g.remain, r, r-1) { return true, "allowed" }
    }
}

// ScoreGate denies peers whose score is below threshold.
type ScoreGate struct{
    threshold int64
    scores map[PeerID]int64
}

func NewScoreGate(threshold int64, scores map[PeerID]int64) ScoreGate {
    if scores == nil { scores = map[PeerID]int64{} }
    return ScoreGate{threshold: threshold, scores: scores}
}

func (g ScoreGate) Allow(id PeerID) bool {
    ok, _ := g.AllowWithReason(id)
    return ok
}

func (g ScoreGate) AllowWithReason(id PeerID) (bool, string) {
    if g.threshold <= 0 { return true, "allowed" }
    s, ok := g.scores[id]
    if !ok || s < g.threshold { return false, "scored_out" }
    return true, "allowed"
}