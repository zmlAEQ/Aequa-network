package qbft

import (
    "testing"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// FuzzBasicVerifier_NoPanic exercises qbft.BasicVerifier with structurally valid
// inputs that may violate semantic constraints (e.g., type-specific rounds/heights),
// ensuring no panics and stable metric/log paths.
func FuzzBasicVerifier_NoPanic(f *testing.F) {
    // Seed a variety of representative cases
    f.Add("p", string(MsgPreprepare), uint64(0), uint64(0), uint8(0), true)   // valid preprepare (round=0)
    f.Add("p", string(MsgPreprepare), uint64(10), uint64(1), uint8(8), true)  // preprepare invalid round, short sig
    f.Add("p", string(MsgPrepare),    uint64(0),  uint64(0), uint8(0), true)  // prepare invalid round
    f.Add("p", string(MsgCommit),     uint64(1),  uint64(1), uint8(64), true) // commit ok
    f.Add("q", "bogus",               uint64(1),  uint64(0), uint8(0), false) // invalid type

    f.Fuzz(func(t *testing.T, from string, typ string, height uint64, round uint64, sigLen uint8, allow bool) {
        metrics.Reset()
        v := NewBasicVerifier()
        // Configure global and type-scoped windows (placeholders) to widen coverage
        v.SetReplayWindow(2)
        v.SetTypeMinHeight(MsgPrepare, 100) // prepare must be at least height 100
        v.SetTypeRoundMax(MsgCommit, 3)     // commit round upper bound = 3
        if allow {
            v.SetAllowed(from)
        } else {
            v.SetAllowed(from+"_other") // likely unauthorized
        }

        // Build a message from fuzzed fields
        sig := make([]byte, int(sigLen)) // len<32 hits sig-shape error path; 0 is allowed
        msg := Message{
            ID:     "fuzz",       // stable id for replay coverage below
            From:   from,
            Type:   Type(typ),     // may be invalid -> should return error (no panic)
            Height: height,
            Round:  round,
            Sig:    sig,
        }

        // Single verification (should never panic)
        _ = v.Verify(msg)

        // Replay with same id (triggers anti-replay path, no panic)
        _ = v.Verify(msg)

        // Move height forward beyond replay window; should not panic
        msg.Height += 3
        _ = v.Verify(msg)
    })
}

