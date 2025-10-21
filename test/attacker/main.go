package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
)

type Message struct {
    From    string `json:"From"`
    Height  uint64 `json:"Height"`
    Round   uint64 `json:"Round"`
    Type    string `json:"Type"`
    Payload []byte `json:"Payload"`
    ID      string `json:"ID"`
    TraceID string `json:"TraceID"`
    Sig     []byte `json:"Sig"`
}

func postJSON(url string, v any) error {
    b, _ := json.Marshal(v)
    resp, err := http.Post(url, "application/json", bytes.NewReader(b))
    if err != nil { return err }
    _ = resp.Body.Close()
    return nil
}

func main() {
    rand.Seed(time.Now().UnixNano())
    eps := os.Getenv("ENDPOINTS")
    if eps == "" { eps = "http://127.0.0.1:4610,http://127.0.0.1:4611,http://127.0.0.1:4612,http://127.0.0.1:4613" }
    targets := strings.Split(eps, ",")
    urls := make([]string, 0, len(targets))
    for _, t := range targets { t = strings.TrimSpace(t); if t != "" { urls = append(urls, fmt.Sprintf("%s/e2e/qbft", t)) } }
    if len(urls) == 0 { fmt.Println("no endpoints"); os.Exit(1) }

    var storedCommit *Message
    ticker := time.NewTicker(1500 * time.Millisecond)
    defer ticker.Stop()
    for {
        <-ticker.C
        u := urls[rand.Intn(len(urls))]
        // Random scenario selector
        sc := rand.Intn(6)
        switch sc {
        case 0: // valid preprepare then store a commit for replay
            h := uint64(rand.Intn(1000) + 1)
            pre := Message{From:"L", Type:"preprepare", Height:h, Round:0, ID:fmt.Sprintf("blk-%d", h)}
            _ = postJSON(u, pre)
            cm := Message{From:"C1", Type:"commit", Height:h, Round:1, ID:pre.ID}
            _ = postJSON(u, cm)
            storedCommit = &cm
        case 1: // replay stored commit
            if storedCommit != nil { _ = postJSON(u, *storedCommit) }
        case 2: // commit with mismatched proposal id (safety)
            h := uint64(rand.Intn(1000) + 1)
            pre := Message{From:"L", Type:"preprepare", Height:h, Round:0, ID:fmt.Sprintf("blk-%d", h)}
            _ = postJSON(u, pre)
            bad := Message{From:"C2", Type:"commit", Height:h, Round:1, ID:"other"}
            _ = postJSON(u, bad)
        case 3: // wrong round semantics
            _ = postJSON(u, Message{From:"P1", Type:"prepare", Height:10, Round:0, ID:"x"})
            _ = postJSON(u, Message{From:"C1", Type:"commit", Height:11, Round:0, ID:"y"})
        case 4: // short signature shape
            _ = postJSON(u, Message{From:"S", Type:"prepare", Height:12, Round:1, ID:"z", Sig:[]byte{1,2,3}})
        default: // ok prepare to keep liveness visible
            h := uint64(rand.Intn(1000) + 1)
            pre := Message{From:"L", Type:"preprepare", Height:h, Round:0, ID:fmt.Sprintf("blk-%d", h)}
            _ = postJSON(u, pre)
            _ = postJSON(u, Message{From:"P1", Type:"prepare", Height:h, Round:1, ID:pre.ID})
            _ = postJSON(u, Message{From:"P2", Type:"prepare", Height:h, Round:1, ID:pre.ID})
        }
    }
}

