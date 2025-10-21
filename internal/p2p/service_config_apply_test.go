package p2p

import (
    "context"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/internal/dkg"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

func TestP2P_Config_AllowList_Applies_First(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(Config{MaxConns: 8, AllowList: []PeerID{"P1"}})
    if err := s.Start(context.Background()); err != nil { t.Fatalf("start: %v", err) }

    if err := s.Connect("P2"); err == nil { t.Fatalf("P2 should be denied by allowlist") }
    if err := s.Connect("P1"); err != nil { t.Fatalf("P1 should be allowed: %v", err) }

    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="denied"} 1`) {
        t.Fatalf("want denied=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="allowed"} 1`) {
        t.Fatalf("want allowed=1, got %q", dump)
    }
}

func TestP2P_Config_RateLimit_Applies_Limited_Label(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(Config{MaxConns: 8, RateLimit: 1})
    if err := s.Start(context.Background()); err != nil { t.Fatalf("start: %v", err) }

    if err := s.Connect("A"); err != nil { t.Fatalf("first connect allowed: %v", err) }
    if err := s.Connect("B"); err == nil { t.Fatalf("second connect should be limited by rate gate") }

    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="limited"} 1`) {
        t.Fatalf("want limited=1 (rate), got %q", dump)
    }
}

func TestP2P_Config_ScoreThreshold_Denies_WhenNoScores(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(Config{MaxConns: 8, ScoreThreshold: 10})
    if err := s.Start(context.Background()); err != nil { t.Fatalf("start: %v", err) }
    if err := s.Connect("S1"); err == nil { t.Fatalf("expect denial due to score threshold with no scores provided") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="denied"} 1`) {
        t.Fatalf("want denied=1 (score), got %q", dump)
    }
}

func TestP2P_Config_MaxConns_Drives_ResourceLimits(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(Config{MaxConns: 1})
    if err := s.Start(context.Background()); err != nil { t.Fatalf("start: %v", err) }
    if err := s.Connect("X"); err != nil { t.Fatalf("first connect should pass: %v", err) }
    if err := s.Connect("Y"); err == nil { t.Fatalf("second connect should be limited by resources") }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_open_total 1`) {
        t.Fatalf("want open_total=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conns_open 1`) {
        t.Fatalf("want conns_open=1, got %q", dump)
    }
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="limited"} 1`) {
        t.Fatalf("want limited=1 (resource), got %q", dump)
    }
}

func TestP2P_Config_Order_DKG_After_Gates(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(Config{MaxConns: 8, AllowList: []PeerID{"P1", "P3"}})
    s.SetDKG(dkg.NewStaticVerifier("P1"))
    if err := s.Start(context.Background()); err != nil { t.Fatalf("start: %v", err) }
    // P3 passes allowlist but fails DKG
    if err := s.Connect("P3"); err == nil { t.Fatalf("P3 should be denied by DKG") }
    // P1 should pass all
    if err := s.Connect("P1"); err != nil { t.Fatalf("P1 should connect: %v", err) }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_conn_attempts_total{result="dkg_denied"} 1`) {
        t.Fatalf("want dkg_denied=1, got %q", dump)
    }
}

