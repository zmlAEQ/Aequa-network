package payload

import (
    "encoding/json"
    "testing"
)

func TestJSONManager_Validate_SizeAndJSON(t *testing.T) {
    m := NewJSONManager(8)
    if err := m.Validate([]byte(`{"a":1}`)); err != nil {
        t.Fatalf("validate ok json: %v", err)
    }
    if err := m.Validate([]byte(`not-json`)); err == nil {
        t.Fatalf("want invalid json error")
    }
    if err := m.Validate([]byte(`{"size":123456}`)); err == nil { // >8 bytes
        t.Fatalf("want size error")
    }
}

func TestJSONManager_Canonical_HashStable(t *testing.T) {
    m := NewJSONManager(1<<20)
    a := []byte(`{"x":1,  "y": [1,2,3]}`)
    b := []byte("\n { \n \t\"x\":1, \"y\":[1,2,3] }  ")
    ca, err := m.Canonical(a); if err != nil { t.Fatalf("canon a: %v", err) }
    cb, err := m.Canonical(b); if err != nil { t.Fatalf("canon b: %v", err) }
    if !jsonEqual(ca, cb) { t.Fatalf("canonical mismatch: %q vs %q", ca, cb) }
    ha, _ := m.Hash(a)
    hb, _ := m.Hash(b)
    if ha != hb { t.Fatalf("hash mismatch for semantically equal payloads") }
}

func jsonEqual(a, b []byte) bool {
    var va, vb any
    if json.Unmarshal(a, &va) != nil { return false }
    if json.Unmarshal(b, &vb) != nil { return false }
    return jsonMarshalString(va) == jsonMarshalString(vb)
}

func jsonMarshalString(v any) string {
    b, _ := json.Marshal(v)
    return string(b)
}