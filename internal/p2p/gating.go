package p2p

// Gate decides whether a peer is allowed to connect.
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