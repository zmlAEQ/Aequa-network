package tss

// Package tss provides a tiny placeholder for threshold-signature style
// aggregation so that the tss-logic-stress workflow has a concrete target to
// fuzz against. The implementation is intentionally simple and self-contained
// (no external deps) and only encodes a few invariants we can fuzz: non-empty
// input, commutativity and idempotency.

// Share is a minimal stand-in for a partial signature share.
type Share uint64

// Aggregate combines a set of shares into a single aggregate value. It returns
// an error when the input is empty. For fuzzability we intentionally keep the
// behaviour deterministic, order-insensitive and side-effect free.
func Aggregate(shares []Share) (Share, error) {
    if len(shares) == 0 {
        return 0, ErrEmpty
    }
    var s uint64
    for _, v := range shares {
        s += uint64(v)
    }
    return Share(s), nil
}

// VerifyAggregate is a no-op placeholder that asserts the aggregate is derived
// from at least one share. Real verification would check against a group
// element or public key; here we only keep a sanity predicate to fuzz.
func VerifyAggregate(agg Share) bool { return agg != 0 }

// ErrEmpty is returned when Aggregate is called with no shares.
type Err string
func (e Err) Error() string { return string(e) }

const ErrEmpty Err = "empty shares"

