package p2p

import "testing"

type countHook struct{ j, l int }
func (h *countHook) OnPeerJoin(id string)  { h.j++ }
func (h *countHook) OnPeerLeave(id string) { h.l++ }

func TestHook_IsInvokedOnConnectDisconnect(t *testing.T) {
    h := &countHook{}
    s := NewWithOpts(nil, AllowAllGate{}, NewResourceManager(ResourceLimits{MaxConns: 1}), h)
    if err := s.Connect("A"); err != nil { t.Fatalf("connect: %v", err) }
    s.Disconnect("A")
    if h.j != 1 || h.l != 1 {
        t.Fatalf("want join=1, leave=1; got %d,%d", h.j, h.l)
    }
}