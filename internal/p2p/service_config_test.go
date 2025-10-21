package p2p

import (
    "context"
    "errors"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// stub verifier failing cluster verification
type failVerifier struct{}
func (failVerifier) VerifyCluster() error { return errors.New("dkg invalid") }
func (failVerifier) AllowPeer(id string) bool { return true }

func TestP2P_Start_ConfigInvalid_FailFast(t *testing.T) {
    metrics.Reset()
    s := New()
    // invalid MaxConns
    s.SetConfig(Config{MaxConns: -1})
    if err := s.Start(context.Background()); err == nil {
        t.Fatalf("want fail-fast on invalid config")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_config_checks_total{result="error"} 1`) {
        t.Fatalf("want config error metric, got %q", dump)
    }
}

func TestP2P_Start_DKGInvalid_FailFast(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(DefaultConfig())
    s.SetDKG(failVerifier{})
    if err := s.Start(context.Background()); err == nil {
        t.Fatalf("want fail-fast on dkg invalid")
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_dkg_cluster_checks_total{result="error"} 1`) {
        t.Fatalf("want dkg error metric, got %q", dump)
    }
}

func TestP2P_Start_OK(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(DefaultConfig())
    if err := s.Start(context.Background()); err != nil {
        t.Fatalf("start ok: %v", err)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_config_checks_total{result="ok"} 1`) {
        t.Fatalf("want config ok metric, got %q", dump)
    }
}