package config

import (
    "encoding/json"
    "os"
)

type Operator struct { Index int `json:"index"`; PeerID string `json:"peer_id"` }

type ClusterLock struct {
    Name      string     `json:"name"`
    Threshold int        `json:"threshold"`
    Operators []Operator `json:"operators"`
}

func LoadClusterLock(path string) (ClusterLock, error) {
    var c ClusterLock
    b, err := os.ReadFile(path)
    if err != nil { return c, err }
    err = json.Unmarshal(b, &c)
    return c, err
}
