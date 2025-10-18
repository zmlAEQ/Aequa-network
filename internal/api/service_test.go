package api

import "testing"

func TestValidateDutyJSON_Good(t *testing.T) {
	good := []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`)
	if err := validateDutyJSON(good); err != nil { t.Fatalf("unexpected: %v", err) }
}

func TestValidateDutyJSON_Bad(t *testing.T) {
	cases := [][]byte{
		nil,
		[]byte(""),
		[]byte("{"),
		[]byte(`{"type":"x"}`),
		[]byte(`{"type":"attester","height":-1}`),
		[]byte(`{"type":"attester","height":1,"round":-1}`),
	}
	for i, b := range cases {
		if err := validateDutyJSON(b); err == nil {
			t.Fatalf("case %d: want error", i)
		}
	}
}
