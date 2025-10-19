package api

import (
    "context"
    "testing"
)

// TestService_StartStop exercises Start and Stop to cover server lifecycle
// without binding to a fixed port (uses :0 for ephemeral port).
func TestService_StartStop(t *testing.T) {
    s := New("127.0.0.1:0", nil, "")
    if err := s.Start(context.Background()); err != nil {
        t.Fatalf("start error: %v", err)
    }
    if err := s.Stop(context.Background()); err != nil {
        t.Fatalf("stop error: %v", err)
    }
}