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


// FuzzState_NoPanic exercises the qbft.State processor with short sequences
// that combine commit preconditions, duplicates and proposal-ID mismatches,
// asserting only that no panics occur and phases remain within known labels.
func FuzzState_NoPanic(f *testing.F) {
    // Seeds toggle: havePreprepare, duplicatePrepare, commitMismatchedID
    f.Add(uint8(0), uint8(0), uint8(0)) // preprepare + no dup + matched commit
    f.Add(uint8(1), uint8(1), uint8(1)) // no preprepare + dup + mismatched commit
    f.Add(uint8(0), uint8(1), uint8(0)) // preprepare + dup + matched commit

    f.Fuzz(func(t *testing.T, a, b, c uint8) {
        st := &State{Leader: "L"}
        // Optional preprepare
        if a%2 == 0 {
            _ = st.Process(Message{ID: "blk", From: "L", Type: MsgPreprepare, Height: 10, Round: 0})
        }
        // First prepare (may be before preprepare to exercise error path)
        _ = st.Process(Message{ID: "blk", From: "P1", Type: MsgPrepare, Height: 10, Round: 1})
        // Optional duplicate or second distinct prepare to reach threshold when preprepared
        if b%2 == 0 { // distinct second
            _ = st.Process(Message{ID: "blk", From: "P2", Type: MsgPrepare, Height: 10, Round: 1})
        } else { // duplicate
            _ = st.Process(Message{ID: "blk", From: "P1", Type: MsgPrepare, Height: 10, Round: 1})
        }
        // Commit with matched or mismatched proposal ID
        id := "blk"
        if c%2 == 1 { id = "blkX" }
        _ = st.Process(Message{ID: id, From: "C1", Type: MsgCommit, Height: 10, Round: 1})

        // Ensure phase remains one of known labels (empty, preprepared, prepared, commit)
        switch st.Phase {
        case "", "preprepared", "prepared", "commit":
            // ok
        default:
            t.Fatalf("unknown phase label: %q", st.Phase)
        }
    })
}
