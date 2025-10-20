package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestBasicVerifier_OK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    msg := Message{ID:"1", From:"p", Type:MsgPrepare}
    if err := v.Verify(msg); err != nil { t.Fatalf("unexpected: %v", err) }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{type="prepare"} 1`) {
        t.Fatalf("want ok count for prepare, got %q", dump)
    }
}

func TestBasicVerifier_Replay(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    msg := Message{ID:"x", From:"p", Type:MsgCommit}
    _ = v.Verify(msg)
    if err := v.Verify(msg); err == nil { t.Fatalf("want replay error") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="replay"} 1`) {
        t.Fatalf("want replay count, got %q", dump)
    }
}

func TestBasicVerifier_Invalid(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    if err := v.Verify(Message{}); err == nil { t.Fatalf("want invalid error") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="error"} 1`) {
        t.Fatalf("want invalid count, got %q", dump)
    }
}