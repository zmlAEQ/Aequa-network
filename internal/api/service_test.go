package api

import "testing"

func TestValidateDutyJSON(t *testing.T) {
	good := []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`)
	if err := validateDutyJSON(good); err != nil { t.Fatalf("unexpected: %v", err) }
	bad := [][]byte{ nil, []byte(""), []byte("{"), []byte("{}"), []byte(`{"type":"x"}`) }
	for _, b := range bad { if err := validateDutyJSON(b); err == nil { t.Fatalf("want error for %q", string(b)) } }
}

func FuzzValidateDutyJSON(f *testing.F) {
	f.Add([]byte(`{"type":"attester","height":1,"round":0,"payload":{}}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		_ = validateDutyJSON(data) // must not panic
	})
}

