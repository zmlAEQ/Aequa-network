package api

import "testing"

// FuzzValidateDutyJSON fuzzes validateDutyJSON with inline seeds.
func FuzzValidateDutyJSON(f *testing.F) {
    // Inline seeds (migrate away from legacy testdata corpus format)
    seeds := [][]byte{
        []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`),
        []byte(`{"type":"proposer","height":0,"round":0,"payload":{}}`),
        []byte(`{"type":"sync","height":2,"round":3,"payload":{}}`),
        []byte(`{}`),
        []byte(`{"type":"x"}`),
        []byte(`{"type":"attester","height":999999999999,"round":0}`),
    }
    for _, s := range seeds { f.Add(s) }

    f.Fuzz(func(t *testing.T, b []byte) {
        _ = validateDutyJSON(b)
    })
}