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
package api

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandleDuty_Success(t *testing.T) {
    s := &Service{onPublish: func(_ context.Context, p []byte) error { return nil }}
    rr := httptest.NewRecorder()
    body := []byte(`{"type":"attester","height":1,"round":0,"payload":{}}`)
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(body))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusAccepted { t.Fatalf("want 202, got %d", rr.Code) }
}

func TestHandleDuty_MethodNotAllowed(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/v1/duty", nil)
    s.handleDuty(rr, req)
    if rr.Code != http.StatusMethodNotAllowed { t.Fatalf("want 405, got %d", rr.Code) }
}

func TestHandleDuty_EmptyBody(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", nil)
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest { t.Fatalf("want 400, got %d", rr.Code) }
}
func TestValidateDutyJSON_OutOfRange(t *testing.T) {
    b := []byte(`{"type":"attester","height":4611686018427387908,"round":0,"payload":{}}`)
    if err := validateDutyJSON(b); err == nil {
        t.Fatalf("want error for out-of-range height")
    }
}

func TestHandleDuty_OversizeBody(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    big := make([]byte, (1<<20)+1)
    for i := range big { big[i] = 'a' }
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(big))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest { t.Fatalf("want 400, got %d", rr.Code) }
}
