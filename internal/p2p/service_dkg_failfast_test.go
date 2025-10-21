package p2p

import (
    "context"
    "errors"
    "strings"
    "testing"

    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// failDKG returns error on cluster verification, allows all peers otherwise.
type failDKG struct{}
func (failDKG) VerifyCluster() error { return errors.New("cluster invalid") }
func (failDKG) AllowPeer(id string) bool { return true }

// okDKG returns success on cluster verification and allows all peers.
type okDKG struct{}
func (okDKG) VerifyCluster() error { return nil }
func (okDKG) AllowPeer(id string) bool { return true }

// Test that DKG precheck fails fast at Start and records consistent metrics.
func TestP2P_Start_DKGFailFast_MetricsConsistency(t *testing.T) {
    metrics.Reset()
    s := New()
    // Config valid so config check should be ok before DKG fails.
    s.SetConfig(DefaultConfig())
    s.SetDKG(failDKG{})

    if err := s.Start(context.Background()); err == nil {
        t.Fatalf("want start fail-fast when DKG precheck fails")
    }

    dump := metrics.DumpProm()
    // Config validated successfully before DKG step
    if !strings.Contains(dump, `p2p_config_checks_total{result="ok"} 1`) {
        t.Fatalf("want p2p_config_checks_total ok=1, got %q", dump)
    }
    // DKG precheck recorded as error
    if !strings.Contains(dump, `p2p_dkg_cluster_checks_total{result="error"} 1`) {
        t.Fatalf("want p2p_dkg_cluster_checks_total error=1, got %q", dump)
    }
    // service_op_ms summary should include one start observation regardless of outcome
    if !strings.Contains(dump, `service_op_ms_count{op="start",service="p2p"}`) {
        t.Fatalf("want service_op_ms_count for p2p start, got %q", dump)
    }
}

// Test that DKG precheck success path records ok metric and start observation.
func TestP2P_Start_DKGPrecheck_OK_Metrics(t *testing.T) {
    metrics.Reset()
    s := New()
    s.SetConfig(DefaultConfig())
    s.SetDKG(okDKG{})
    if err := s.Start(context.Background()); err != nil {
        t.Fatalf("start ok: %v", err)
    }
    dump := metrics.DumpProm()
    if !strings.Contains(dump, `p2p_dkg_cluster_checks_total{result="ok"} 1`) {
        t.Fatalf("want p2p_dkg_cluster_checks_total ok=1, got %q", dump)
    }
    if !strings.Contains(dump, `service_op_ms_count{op="start",service="p2p"}`) {
        t.Fatalf("want service_op_ms_count for p2p start, got %q", dump)
    }
}

