//go:build e2e

package p2p

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/logger"
)

type p2pReq struct{ ID string `json:"id"` }

// startP2PE2E exposes /e2e/p2p/connect|disconnect on 0.0.0.0:4615 for e2e testing.
func startP2PE2E(s *Service) {
    mux := http.NewServeMux()
    mux.HandleFunc("/e2e/p2p/connect", func(w http.ResponseWriter, r *http.Request) {
        begin := time.Now(); defer r.Body.Close()
        var req p2pReq; _ = json.NewDecoder(r.Body).Decode(&req)
        if req.ID == "" { http.Error(w, "missing id", 400); return }
        err := s.Connect(PeerID(req.ID))
        res := "ok"; if err != nil { res = "denied" }
        logger.InfoJ("p2p_attack", map[string]any{"op":"connect", "peer_id": req.ID, "result": res, "latency_ms": time.Since(begin).Milliseconds()})
        if err != nil { w.WriteHeader(202) } else { w.WriteHeader(200) }
    })
    mux.HandleFunc("/e2e/p2p/disconnect", func(w http.ResponseWriter, r *http.Request) {
        begin := time.Now(); defer r.Body.Close()
        var req p2pReq; _ = json.NewDecoder(r.Body).Decode(&req)
        if req.ID == "" { http.Error(w, "missing id", 400); return }
        s.Disconnect(PeerID(req.ID))
        logger.InfoJ("p2p_attack", map[string]any{"op":"disconnect", "peer_id": req.ID, "result": "ok", "latency_ms": time.Since(begin).Milliseconds()})
        w.WriteHeader(200)
    })
    srv := &http.Server{Addr: "0.0.0.0:4615", Handler: mux}
    go func() { _ = srv.ListenAndServe() }()
}

