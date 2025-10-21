package main

import (
    "bytes"
    "crypto/rand"
    "encoding/binary"
    "encoding/json"
    "fmt"
    mrand "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
)

// qbft message used by /e2e/qbft endpoint
type qbftMsg struct {
    From    string `json:"From"`
    Height  uint64 `json:"Height"`
    Round   uint64 `json:"Round"`
    Type    string `json:"Type"`
    Payload []byte `json:"Payload"`
    ID      string `json:"ID"`
    TraceID string `json:"TraceID"`
    Sig     []byte `json:"Sig"`
}

// p2p request used by /e2e/p2p/connect|disconnect endpoints
type p2pReq struct {
    ID string `json:"id"`
}

func postJSON(url string, v any) error {
    b, _ := json.Marshal(v)
    resp, err := http.Post(url, "application/json", bytes.NewReader(b))
    if err != nil { return err }
    _ = resp.Body.Close()
    return nil
}

func main() {
    // local RNG seeded from crypto/rand to avoid deprecated global Seed
    var seed int64
    var b8 [8]byte
    if _, err := rand.Read(b8[:]); err == nil { seed = int64(binary.LittleEndian.Uint64(b8[:])) } else { seed = time.Now().UnixNano() }
    r := mrand.New(mrand.NewSource(seed))

    // endpoints
    epsQBFT := os.Getenv("ENDPOINTS_QBFT")
    if epsQBFT == "" { epsQBFT = "http://127.0.0.1:4610,http://127.0.0.1:4611,http://127.0.0.1:4612,http://127.0.0.1:4613" }
    qbftTargets := splitNonEmpty(epsQBFT)
    for i := range qbftTargets { qbftTargets[i] = strings.TrimRight(qbftTargets[i], "/") + "/e2e/qbft" }

    epsP2P := os.Getenv("ENDPOINTS_P2P")
    if epsP2P == "" { epsP2P = "http://127.0.0.1:4615,http://127.0.0.1:4616,http://127.0.0.1:4617,http://127.0.0.1:4618" }
    p2pTargets := splitNonEmpty(epsP2P)
    for i := range p2pTargets { p2pTargets[i] = strings.TrimRight(p2pTargets[i], "/") }

    // p2p storm parameters
    idSpace := os.Getenv("P2P_ID_PREFIX")
    if idSpace == "" { idSpace = "X" }
    maxIDs := 1024
    if v := os.Getenv("P2P_MAX_IDS"); v != "" { if n, err := fmt.Sscan(v, &maxIDs); err == nil && n == 1 && maxIDs > 0 { /* ok */ } }

    // Start goroutine: qbft adversarial injection
    go func() {
        var stored *qbftMsg
        ticker := time.NewTicker(1500 * time.Millisecond); defer ticker.Stop()
        for range ticker.C {
            u := qbftTargets[r.Intn(len(qbftTargets))]
            sc := r.Intn(6)
            switch sc {
            case 0: // valid preprepare + store commit
                h := uint64(r.Intn(1000) + 1)
                pre := qbftMsg{From:"L", Type:"preprepare", Height:h, Round:0, ID:fmt.Sprintf("blk-%d", h)}
                _ = postJSON(u, pre)
                cm := qbftMsg{From:"C1", Type:"commit", Height:h, Round:1, ID:pre.ID}
                _ = postJSON(u, cm)
                stored = &cm
            case 1: // replay stored
                if stored != nil { _ = postJSON(u, *stored) }
            case 2: // commit id mismatch
                h := uint64(r.Intn(1000) + 1)
                pre := qbftMsg{From:"L", Type:"preprepare", Height:h, Round:0, ID:fmt.Sprintf("blk-%d", h)}
                _ = postJSON(u, pre)
                bad := qbftMsg{From:"C2", Type:"commit", Height:h, Round:1, ID:"other"}
                _ = postJSON(u, bad)
            case 3: // round semantics errors
                _ = postJSON(u, qbftMsg{From:"P1", Type:"prepare", Height:10, Round:0, ID:"x"})
                _ = postJSON(u, qbftMsg{From:"C1", Type:"commit", Height:11, Round:0, ID:"y"})
            case 4: // short signature
                _ = postJSON(u, qbftMsg{From:"S", Type:"prepare", Height:12, Round:1, ID:"z", Sig:[]byte{1,2,3}})
            default: // liveness: prepare quorum
                h := uint64(r.Intn(1000) + 1)
                pre := qbftMsg{From:"L", Type:"preprepare", Height:h, Round:0, ID:fmt.Sprintf("blk-%d", h)}
                _ = postJSON(u, pre)
                _ = postJSON(u, qbftMsg{From:"P1", Type:"prepare", Height:h, Round:1, ID:pre.ID})
                _ = postJSON(u, qbftMsg{From:"P2", Type:"prepare", Height:h, Round:1, ID:pre.ID})
            }
        }
    }()

    // Start goroutine: p2p connect/disconnect storm
    go func() {
        ticker := time.NewTicker(500 * time.Millisecond); defer ticker.Stop()
        counter := 0
        for range ticker.C {
            base := p2pTargets[r.Intn(len(p2pTargets))]
            // 80% connect, 20% disconnect
            id := fmt.Sprintf("%s-%d", idSpace, r.Intn(maxIDs))
            if r.Intn(5) > 0 {
                _ = postJSON(base+"/e2e/p2p/connect", p2pReq{ID:id})
            } else {
                _ = postJSON(base+"/e2e/p2p/disconnect", p2pReq{ID:id})
            }
            counter++
            if counter%200 == 0 { time.Sleep(2 * time.Second) }
        }
    }()

    select {}
}

func splitNonEmpty(s string) []string {
    parts := strings.Split(s, ",")
    out := make([]string, 0, len(parts))
    for _, p := range parts { p = strings.TrimSpace(p); if p != "" { out = append(out, p) } }
    return out
}

