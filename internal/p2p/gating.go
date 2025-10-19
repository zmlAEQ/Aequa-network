package p2p

// Gate decides whether a peer is allowed to connect.
type Gate interface { Allow(id PeerID) bool }

// AllowAllGate admits all peers (default in minimal setup).
type AllowAllGate struct{}
func (AllowAllGate) Allow(id PeerID) bool { return true }

// StaticDenyGate denies all peers (used in tests).
type StaticDenyGate struct{}
func (StaticDenyGate) Allow(id PeerID) bool { return false }