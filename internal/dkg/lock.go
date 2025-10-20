package dkg

// Verifier exposes methods to validate cluster-lock and peer admission.
type Verifier interface {
    VerifyCluster() error
    AllowPeer(id string) bool
}

// NoopVerifier allows everything and reports success.
type NoopVerifier struct{}

func (NoopVerifier) VerifyCluster() error { return nil }
func (NoopVerifier) AllowPeer(id string) bool { return true }

// StaticVerifier is a simple in-memory verifier based on an allowlist.
type StaticVerifier struct{ allowed map[string]struct{} }

func NewStaticVerifier(ids ...string) StaticVerifier {
    m := make(map[string]struct{}, len(ids))
    for _, id := range ids { m[id] = struct{}{} }
    return StaticVerifier{allowed: m}
}

func (v StaticVerifier) VerifyCluster() error { return nil }
func (v StaticVerifier) AllowPeer(id string) bool { _, ok := v.allowed[id]; return ok }