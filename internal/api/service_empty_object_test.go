package api

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Empty object should be rejected due to missing required fields.
func TestHandleDuty_EmptyObject(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader([]byte(`{}`)))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("want 400, got %d", rr.Code)
    }
}

// Unknown fields should be ignored by json.Unmarshal, but missing required
// fields must still cause a validation error.
func TestHandleDuty_UnknownFieldsOnly(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader([]byte(`{"foo":1,"bar":"x"}`)))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("want 400, got %d", rr.Code)
    }
}

// Boundary values equal to the allowed maximum should pass validation.
func TestHandleDuty_BoundaryAccepted(t *testing.T) {
    s := &Service{}
    rr := httptest.NewRecorder()
    // height == 1<<62 and round == 1<<40 are within allowed range (strict '>' check)
    body := []byte(`{"type":"proposer","height":4611686018427387904,"round":1099511627776,"payload":{}}`)
    req := httptest.NewRequest(http.MethodPost, "/v1/duty", bytes.NewReader(body))
    s.handleDuty(rr, req)
    if rr.Code != http.StatusAccepted {
        t.Fatalf("want 202 for boundary values, got %d", rr.Code)
    }
}