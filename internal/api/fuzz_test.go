package api

import "testing"

// FuzzValidateDutyJSON uses go fuzzing with existing seeds in testdata/.
func FuzzValidateDutyJSON(f *testing.F) {
    f.Fuzz(func(t *testing.T, b []byte) {
        _ = validateDutyJSON(b)
    })
}