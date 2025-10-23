package tss

import (
    "testing"
)

// FuzzAggregate exercises the minimal invariants of the placeholder
// aggregation: non-empty input, commutativity (order-insensitive) and
// idempotency (aggregating the same set twice yields same result).
func FuzzAggregate(f *testing.F) {
    f.Add(uint64(1), uint64(2), uint64(3))
    f.Add(uint64(0), uint64(0), uint64(1))
    f.Fuzz(func(t *testing.T, a, b, c uint64) {
        shares := []Share{Share(a), Share(b), Share(c)}

        // Build a non-empty slice by ensuring at least one element.
        nonEmpty := shares
        if len(nonEmpty) == 0 {
            nonEmpty = []Share{1}
        }

        agg1, err := Aggregate(nonEmpty)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if !VerifyAggregate(agg1) {
            t.Fatalf("aggregate failed VerifyAggregate")
        }

        // Commutativity: swap order should not change result.
        swapped := []Share{nonEmpty[len(nonEmpty)-1]}
        if len(nonEmpty) > 1 {
            swapped = append(swapped, nonEmpty[:len(nonEmpty)-1]...)
        }
        agg2, err := Aggregate(swapped)
        if err != nil {
            t.Fatalf("unexpected error on swapped: %v", err)
        }
        if agg1 != agg2 {
            t.Fatalf("non-commutative: %v vs %v", agg1, agg2)
        }

        // Idempotency on the same set again.
        agg3, _ := Aggregate(nonEmpty)
        if agg1 != agg3 {
            t.Fatalf("non-idempotent: %v vs %v", agg1, agg3)
        }
    })
}

