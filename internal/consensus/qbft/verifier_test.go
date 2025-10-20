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

func TestBasicVerifier_TypeMinHeight_Old(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetTypeMinHeight(MsgPrepare, 100)
    // below type-scoped min should be rejected
    if err := v.Verify(Message{ID:"tmin-old", From:"p", Type:MsgPrepare, Height:99, Round:1}); err == nil {
        t.Fatalf("want type-scoped old")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="old"} 1`) {
        t.Fatalf("want old=1, got %q", dump)
    }
}

func TestBasicVerifier_TypeMinHeight_BoundaryOK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetTypeMinHeight(MsgPrepare, 100)
    if err := v.Verify(Message{ID:"tmin-ok", From:"p", Type:MsgPrepare, Height:100, Round:1}); err != nil {
        t.Fatalf("boundary should pass: %v", err)
    }
}

func TestBasicVerifier_TypeRoundMax_OOB(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetTypeRoundMax(MsgCommit, 3)
    if err := v.Verify(Message{ID:"tmax-oob", From:"p", Type:MsgCommit, Round:4}); err == nil {
        t.Fatalf("want type-scoped round_oob")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="round_oob"} 1`) {
        t.Fatalf("want round_oob=1, got %q", dump)
    }
}

func TestBasicVerifier_TypeRoundMax_BoundaryOK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetTypeRoundMax(MsgCommit, 3)
    if err := v.Verify(Message{ID:"tmax-ok", From:"p", Type:MsgCommit, Round:3}); err != nil {
        t.Fatalf("boundary should pass: %v", err)
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

// ---- Priority and precedence tests ----

// Global minHeight vs type-scoped minHeight: stricter threshold must prevail.
func TestBasicVerifier_MinHeight_Priority_TypeStricter(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    v.SetMinHeight(50)
    v.SetTypeMinHeight(MsgPrepare, 100)
    // height 80 is >= global 50 but < type 100: should be old due to type rule
    if err := v.Verify(Message{ID: "prio-h-type", From: "p", Type: MsgPrepare, Height: 80, Round: 1}); err == nil {
        t.Fatalf("want old height due to stricter type minHeight")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="old"} 1`) {
        t.Fatalf("want old=1, got %q", dump)
    }
}

func TestBasicVerifier_MinHeight_Priority_GlobalStricter(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    v.SetMinHeight(120)
    v.SetTypeMinHeight(MsgPrepare, 100)
    // height 110 is >= type 100 but < global 120: should be old due to global rule
    if err := v.Verify(Message{ID: "prio-h-global", From: "p", Type: MsgPrepare, Height: 110, Round: 1}); err == nil {
        t.Fatalf("want old height due to stricter global minHeight")
    }
}

// Global roundWindow vs type-scoped roundMax: stricter cap must prevail.
func TestBasicVerifier_RoundWindow_Priority_TypeStricter(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    v.SetRoundWindow(5)
    v.SetTypeRoundMax(MsgCommit, 3)
    if err := v.Verify(Message{ID: "prio-r-type", From: "p", Type: MsgCommit, Round: 4}); err == nil {
        t.Fatalf("want type-scoped round_oob when type is stricter")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="round_oob"} 1`) {
        t.Fatalf("want round_oob=1, got %q", dump)
    }
}

func TestBasicVerifier_RoundWindow_Priority_GlobalStricter(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier()
    v.SetRoundWindow(3)
    v.SetTypeRoundMax(MsgCommit, 5)
    if err := v.Verify(Message{ID: "prio-r-global", From: "p", Type: MsgCommit, Round: 4}); err == nil {
        t.Fatalf("want global round_oob when global is stricter")
    }
}

// Replay precedence: when replayWindow>0, use height-windowed; when ==0, use id-level.
func TestBasicVerifier_Replay_Preference_Windowed(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetReplayWindow(2)
    // first ok
    if err := v.Verify(Message{ID: "rw-pref", From: "p", Type: MsgPrepare, Height: 100, Round: 1}); err != nil {
        t.Fatalf("first msg: %v", err)
    }
    // within window → replay
    if err := v.Verify(Message{ID: "rw-pref", From: "p", Type: MsgPrepare, Height: 101, Round: 1}); err == nil {
        t.Fatalf("want replay within window")
    }
    // outside window → ok
    if err := v.Verify(Message{ID: "rw-pref", From: "p", Type: MsgPrepare, Height: 103, Round: 1}); err != nil {
        t.Fatalf("outside window should pass: %v", err)
    }
}

func TestBasicVerifier_Replay_Preference_IdLevel(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetReplayWindow(0)
    if err := v.Verify(Message{ID: "id-pref", From: "p", Type: MsgCommit, Round: 1}); err != nil {
        t.Fatalf("first commit should pass: %v", err)
    }
    if err := v.Verify(Message{ID: "id-pref", From: "p", Type: MsgCommit, Round: 1}); err == nil {
        t.Fatalf("want id-level replay on second commit")
    }
}

func TestBasicVerifier_Commit_OldHeightReject(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetMinHeight(10)
    if err := v.Verify(Message{ID:"hc9", From:"p", Type:MsgCommit, Height:9, Round:1}); err == nil {
        t.Fatalf("want old height error for commit")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="old"} 1`) {
        t.Fatalf("want old=1, got %q", dump)
    }
}

func TestBasicVerifier_Commit_MinHeightBoundaryOK(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetMinHeight(10)
    if err := v.Verify(Message{ID:"hc10", From:"p", Type:MsgCommit, Height:10, Round:1}); err != nil {
        t.Fatalf("boundary height should pass for commit: %v", err)
    }
}

func TestBasicVerifier_Commit_TypeMinHeight_Old(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetTypeMinHeight(MsgCommit, 100)
    if err := v.Verify(Message{ID:"tc-old", From:"p", Type:MsgCommit, Height:99, Round:1}); err == nil {
        t.Fatalf("want type-scoped old for commit")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="old"} 1`) {
        t.Fatalf("want old=1, got %q", dump)
    }
}

func TestBasicVerifier_Commit_SenderUnauthorized(t *testing.T) {
    metrics.Reset()
    v := NewBasicVerifier(); v.SetAllowed("p")
    if err := v.Verify(Message{ID:"uc1", From:"q", Type:MsgCommit, Round:1}); err == nil {
        t.Fatalf("want unauthorized for commit from q")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `qbft_msg_verified_total{result="unauthorized"} 1`) {
        t.Fatalf("want unauthorized=1, got %q", dump)
    }
}
