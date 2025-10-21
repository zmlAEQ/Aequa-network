//go:build e2e

package consensus

import (
    "encoding/json"
    "net/http"
    "time"

    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
)

// startE2E launches a minimal HTTP server (0.0.0.0:4610) exposing /e2e/qbft
// to inject qbft.Message directly into verifier/state for adversarial testing.
// Compiled only in builds with -tags e2e. Production builds include a no-op.
func startE2E(s *Service) {
    mux := http.NewServeMux()
    mux.HandleFunc("/e2e/qbft", func(w http.ResponseWriter, r *http.Request) {
        begin := time.Now()
        defer r.Body.Close()
        var msg qbft.Message
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
            http.Error(w, "bad json", http.StatusBadRequest)
            return
        }
        status := http.StatusOK
        result := "ok"
        if err := s.v.Verify(msg); err != nil {
            result = "rejected"
            status = http.StatusAccepted // treated as observed but not processed
        } else {
            _ = s.st.Process(msg)
        }
        logger.InfoJ("qbft_attack", map[string]any{
            "result":    result,
            "type":      string(msg.Type),
            "height":    msg.Height,
            "round":     msg.Round,
            "id":        msg.ID,
            "trace_id":  msg.TraceID,
            "latency_ms": time.Since(begin).Milliseconds(),
        })
        w.WriteHeader(status)
    })
    srv := &http.Server{Addr: "0.0.0.0:4610", Handler: mux}
    go func() { _ = srv.ListenAndServe() }()
}

