package api

import (
    "strings"
    "testing"
)

// FuzzValidateDutyJSON uses go fuzzing with existing seeds.
func FuzzValidateDutyJSON(f *testing.F) {
    // Inline seeds (migrate away from legacy testdata corpus format)
    seeds := [][]byte{
        []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`),
        []byte(`{"type":"proposer","height":0,"round":0,"payload":{}}`),
        []byte(`{"type":"sync","height":2,"round":3,"payload":{}}`),
        []byte(`{}`),
        []byte(`{"type":"x"}`),
        []byte(`{"type":"attester","height":999999999999,"round":0}`),
        {0xff, 0xfe, 0xfd},
        []byte(`{"type":"attester","height":0,"round":0,"payload":{"nested":{"k":"v"}}}`),
        []byte(strings.Repeat("a", 1024)),
    }
    for _, s := range seeds { f.Add(s) }

    f.Fuzz(func(t *testing.T, b []byte) {
        _ = validateDutyJSON(b)
    })
}