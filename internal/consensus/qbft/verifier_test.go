package qbft

import (
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestBasicVerifier_OK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    msg := Message{ID:"1", From:"p", Type:MsgPrepare, Round:1}
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

func TestBasicVerifier_Preprepare_RoundZero_OK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    msg := Message{ID:"pp0", From:"p", Type:MsgPreprepare, Round:0}
    if err := v.Verify(msg); err != nil { t.Fatalf("unexpected: %v", err) }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{type="preprepare"} 1`) {
        t.Fatalf("want ok count for preprepare, got %q", dump)
    }
}

func TestBasicVerifier_Preprepare_RoundNonZero_Error(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    msg := Message{ID:"pp1", From:"p", Type:MsgPreprepare, Round:1}
    if err := v.Verify(msg); err == nil { t.Fatalf("want preprepare round semantic error") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="error"} 1`) {
        t.Fatalf("want error=1, got %q", dump)
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
func TestBasicVerifier_OldHeightReject(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    v.SetMinHeight(10)
    msg := Message{ID:"h9", From:"p", Type:MsgPrepare, Height:9, Round:1}
    if err := v.Verify(msg); err == nil {
        t.Fatalf("want old height error")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="old"} 1`) {
        t.Fatalf("want old=1, got %q", dump)
    }
}

func TestBasicVerifier_InvalidType(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    bad := Message{ID:"z", From:"p", Type:Type("bogus")}
    if err := v.Verify(bad); err == nil {
        t.Fatalf("want invalid type error")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="error"} 1`) {
        t.Fatalf("want error=1, got %q", dump)
    }
}
func TestBasicVerifier_RoundWindow(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetRoundWindow(5)
    msg := Message{ID:"rw", From:"p", Type:MsgCommit, Round:6}
    if err := v.Verify(msg); err == nil { t.Fatalf("want round_oob") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="round_oob"} 1`) {
        t.Fatalf("want round_oob=1, got %q", dump)
    }
}

func TestBasicVerifier_SenderUnauthorized(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetAllowed("p")
    if err := v.Verify(Message{ID:"u1", From:"q", Type:MsgPrepare}); err == nil { t.Fatalf("want unauthorized") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="unauthorized"} 1`) {
        t.Fatalf("want unauthorized=1, got %q", dump)
    }
}

func TestBasicVerifier_SignatureShape(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    if err := v.Verify(Message{ID:"s1", From:"p", Type:MsgPrepare, Sig: make([]byte, 8)}); err == nil { t.Fatalf("want sig_invalid") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="sig_invalid"} 1`) {
        t.Fatalf("want sig_invalid=1, got %q", dump)
    }
}

func TestBasicVerifier_MinHeightBoundaryAccept(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetMinHeight(10)
    // boundary height == minHeight should be accepted
    if err := v.Verify(Message{ID:"hb", From:"p", Type:MsgPrepare, Height:10, Round:1}); err != nil {
        t.Fatalf("boundary height should pass: %v", err)
    }
}

func TestBasicVerifier_ReplayWithinWindow(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetReplayWindow(2)
    // same id at height 100 then 101 within window
    if err := v.Verify(Message{ID:"rid", From:"p", Type:MsgPrepare, Height:100, Round:1}); err != nil { t.Fatalf("first msg: %v", err) }
    if err := v.Verify(Message{ID:"rid", From:"p", Type:MsgPrepare, Height:101, Round:1}); err == nil { t.Fatalf("want replay within window") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="replay"} 1`) {
        t.Fatalf("want replay count, got %q", dump)
    }
}

func TestBasicVerifier_ReplayOutsideWindowAccept(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetReplayWindow(2)
    if err := v.Verify(Message{ID:"rid2", From:"p", Type:MsgPrepare, Height:100, Round:1}); err != nil { t.Fatalf("first msg: %v", err) }
    if err := v.Verify(Message{ID:"rid2", From:"p", Type:MsgPrepare, Height:103, Round:1}); err != nil { t.Fatalf("outside window should pass: %v", err) }
}

func TestBasicVerifier_Prepare_RoundZero_Error(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    if err := v.Verify(Message{ID:"pr0", From:"p", Type:MsgPrepare, Round:0}); err == nil {
        t.Fatalf("want prepare round semantic error")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="error"} 1`) {
        t.Fatalf("want error=1, got %q", dump)
    }
}

func TestBasicVerifier_Commit_RoundZero_Error(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    if err := v.Verify(Message{ID:"cm0", From:"p", Type:MsgCommit, Round:0}); err == nil {
        t.Fatalf("want commit round semantic error")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="error"} 1`) {
        t.Fatalf("want error=1, got %q", dump)
    }
}

func TestBasicVerifier_Commit_RoundOne_OK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    if err := v.Verify(Message{ID:"cm1", From:"p", Type:MsgCommit, Round:1}); err != nil {
        t.Fatalf("commit round 1 should pass: %v", err)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{type="commit"} 1`) {
        t.Fatalf("want ok count for commit, got %q", dump)
    }
}
