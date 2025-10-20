package qbft

type Type string

const (
    MsgPreprepare Type = "preprepare"
    MsgPrepare   Type = "prepare"
    MsgCommit    Type = "commit"
)

type Message struct {
    From    string
    Height  uint64
    Round   uint64
    Type    Type
    Payload []byte
    ID      string
    TraceID string
}