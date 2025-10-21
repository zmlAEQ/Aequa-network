//go:build e2e

package p2p

import (
    "fmt"
    "os"
    "strconv"
)

// failVerifier implements dkg.Verifier with failing VerifyCluster.
type failVerifier struct{}

func (failVerifier) VerifyCluster() error { return fmt.Errorf("e2e: dkg invalid") }
func (failVerifier) AllowPeer(id string) bool { return true }

// applyE2E applies e2e-only knobs via env: E2E_P2P_MAXCONNS, E2E_DKG_MODE.
func applyE2E(s *Service) {
    if v := os.Getenv("E2E_P2P_MAXCONNS"); v != "" {
        if n, err := strconv.ParseInt(v, 10, 64); err == nil && n >= 0 {
            s.cfg.MaxConns = n
        }
    }
    if os.Getenv("E2E_DKG_MODE") == "invalid" {
        s.SetDKG(failVerifier{})
    }
}
