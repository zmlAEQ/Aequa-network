package api

import (
    "testing"
)

// Additional fuzz seeds to widen input space coverage.
func FuzzValidateDutyJSON_MoreSeeds(f *testing.F) {
    extra := [][]byte{
        // Nested arrays and objects
        []byte(`{"type":"attester","height":0,"round":0,"payload":{"arr":[1,2,3],"obj":{"k":"v"}}}`),
        // Very long string value in payload
        []byte(`{"type":"sync","height":1,"round":1,"payload":{"s":"` + string(make([]byte, 256)) + `"}}`),
        // Unknown top-level fields alongside valid ones
        []byte(`{"type":"proposer","height":2,"round":0,"unknown":true,"payload":{}}`),
        // Truncated UTF-8 sequence
        {0xe2, 0x82},
        // Random punctuation
        []byte("[{},,,,]"),
    }
    for _, s := range extra { f.Add(s) }

    f.Fuzz(func(t *testing.T, b []byte) {
        _ = validateDutyJSON(b)
    })
}