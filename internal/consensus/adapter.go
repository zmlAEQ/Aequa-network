package consensus

import (
    "fmt"

    qbft "github.com/zmlAEQ/Aequa-network/internal/consensus/qbft"
    "github.com/zmlAEQ/Aequa-network/pkg/bus"
)

// MapEventToQBFT converts a bus.Event into a qbft.Message.
// Stub mapping: direct field mapping with conservative defaults.
func MapEventToQBFT(ev bus.Event) qbft.Message {
    id := fmt.Sprintf("ev-%s-%d-%d", ev.TraceID, ev.Height, ev.Round)
    return qbft.Message{
        ID:      id,
        From:    "consensus_stub",
        Type:    qbft.MsgPrepare,
        Height:  ev.Height,
        Round:   ev.Round,
        Payload: nil,
        TraceID: ev.TraceID,
        Sig:     nil,
    }
}

