package payload

import (
    "crypto/sha256"
    "encoding/json"
    "errors"
)

// Manager defines a minimal interface to validate and canonicalise
// payload bytes without changing existing metrics/logging dimensions.
// The goal in M3 is to provide a safe placeholder that can be unit tested
// and injected later without altering current code paths.
type Manager interface {
    Validate(b []byte) error
    Canonical(b []byte) ([]byte, error)
    Hash(b []byte) ([32]byte, error)
    MaxSize() int
}

var (
    errTooLarge    = errors.New("payload too large")
    errInvalidJSON = errors.New("invalid json payload")
)

// JSONManager performs basic size and JSON structural validation and produces
// a canonical JSON representation using encoding/json. Map key order is not
// strictly defined by the standard library; Hash() therefore uses the canonical
// bytes returned by Canonical to achieve input-normalisation first.
type JSONManager struct {
    max int
}

// NewJSONManager returns a JSONManager with a given size limit in bytes.
func NewJSONManager(maxBytes int) *JSONManager { return &JSONManager{max: maxBytes} }

func (m *JSONManager) MaxSize() int { return m.max }

// Validate enforces size and JSON structural validity.
func (m *JSONManager) Validate(b []byte) error {
    if m.max > 0 && len(b) > m.max {
        return errTooLarge
    }
    var v any
    if err := json.Unmarshal(b, &v); err != nil {
        return errInvalidJSON
    }
    return nil
}

// Canonical returns a compact JSON encoding (no extra whitespace). It performs
// a round-trip through json.Unmarshal to normalise insignificant differences.
func (m *JSONManager) Canonical(b []byte) ([]byte, error) {
    if err := m.Validate(b); err != nil {
        return nil, err
    }
    var v any
    _ = json.Unmarshal(b, &v) // safe after Validate
    enc, err := json.Marshal(v)
    if err != nil {
        return nil, errInvalidJSON
    }
    return enc, nil
}

// Hash returns sha256 over the canonical JSON bytes.
func (m *JSONManager) Hash(b []byte) ([32]byte, error) {
    var zero [32]byte
    enc, err := m.Canonical(b)
    if err != nil { return zero, err }
    return sha256.Sum256(enc), nil
}