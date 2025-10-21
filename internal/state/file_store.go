package state

import (
    "context"
    "encoding/binary"
    "errors"
    "hash/crc32"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// FileStore is a minimal durable Store that persists the last state
// to a single file using an atomic write protocol (tmp write + fsync + rename)
// and a best-effort pessimistic recovery using a backup file.
type FileStore struct {
    mu   sync.Mutex
    path string // main path, e.g. laststate.dat
}

// NewFileStore constructs a file-backed store at the provided path.
func NewFileStore(path string) *FileStore { return &FileStore{path: path} }

const (
    magic   uint32 = 0x53544442 // 'STDB' (State-DB)
    version uint16 = 1
)

// on-disk layout:
// [magic u32][version u16][reserved u16][length u32][crc32 u32][payload bytes...]
// payload = Height u64 | Round u64 (big endian)

func writeFileAtomic(path string, s LastState) error {
    dir := filepath.Dir(path)
    tmp := path + ".tmp"

    f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
    if err != nil { return err }
    // Build payload
    var payload [16]byte
    binary.BigEndian.PutUint64(payload[0:8], s.Height)
    binary.BigEndian.PutUint64(payload[8:16], s.Round)
    // Header (with length + crc)
    length := uint32(len(payload))
    crc := crc32.ChecksumIEEE(payload[:])
    var hdr [4 + 2 + 2 + 4 + 4]byte
    off := 0
    binary.BigEndian.PutUint32(hdr[off:], magic); off += 4
    binary.BigEndian.PutUint16(hdr[off:], version); off += 2
    binary.BigEndian.PutUint16(hdr[off:], 0); off += 2 // reserved
    binary.BigEndian.PutUint32(hdr[off:], length); off += 4
    binary.BigEndian.PutUint32(hdr[off:], crc)

    if _, err = f.Write(hdr[:]); err != nil { _ = f.Close(); return err }
    if _, err = f.Write(payload[:]); err != nil { _ = f.Close(); return err }
    if err = f.Sync(); err != nil { _ = f.Close(); return err }
    if err = f.Close(); err != nil { return err }

    // Ensure directory entry update is durable on some filesystems.
    // (Best effort; ignore error if dir can't be opened on platform.)
    if d, err2 := os.Open(dir); err2 == nil { _ = d.Sync(); _ = d.Close() }

    // Keep previous main as backup (best effort)
    bak := path + ".bak"
    if _, err := os.Stat(path); err == nil {
        // rename current main to .bak; ignore error to avoid blocking progress
        _ = os.Rename(path, bak)
    }
    // Atomic promote tmp -> main
    if err = os.Rename(tmp, path); err != nil { return err }
    if d, err2 := os.Open(dir); err2 == nil { _ = d.Sync(); _ = d.Close() }
    return nil
}

func readFile(path string) (LastState, error) {
    f, err := os.Open(path)
    if err != nil { return LastState{}, err }
    defer f.Close()
    var hdr [4 + 2 + 2 + 4 + 4]byte
    if _, err = io.ReadFull(f, hdr[:]); err != nil { return LastState{}, err }
    off := 0
    mg := binary.BigEndian.Uint32(hdr[off:]); off += 4
    if mg != magic { return LastState{}, errors.New("bad magic") }
    ver := binary.BigEndian.Uint16(hdr[off:]); off += 2
    _ = ver // reserved for future use
    off += 2 // reserved
    length := binary.BigEndian.Uint32(hdr[off:]); off += 4
    wantCRC := binary.BigEndian.Uint32(hdr[off:])
    if length != 16 { return LastState{}, errors.New("bad length") }
    var payload [16]byte
    if _, err = io.ReadFull(f, payload[:]); err != nil { return LastState{}, err }
    got := crc32.ChecksumIEEE(payload[:])
    if got != wantCRC { return LastState{}, errors.New("crc mismatch") }
    s := LastState{
        Height: binary.BigEndian.Uint64(payload[0:8]),
        Round:  binary.BigEndian.Uint64(payload[8:16]),
    }
    return s, nil
}

// SaveLastState persists the last state using atomic file replace.
func (fs *FileStore) SaveLastState(_ context.Context, s LastState) error {
    start := time.Now()
    fs.mu.Lock()
    defer fs.mu.Unlock()
    err := writeFileAtomic(fs.path, s)
    ms := float64(time.Since(start).Milliseconds())
    if err != nil {
        metrics.Inc("state_persist_errors_total", nil)
        logger.ErrorJ("consensus_state", map[string]any{"op":"persist", "result":"error", "err": err.Error()})
        return err
    }
    metrics.ObserveSummary("state_persist_ms", nil, ms)
    logger.InfoJ("consensus_state", map[string]any{"op":"persist", "result":"ok", "height": s.Height, "round": s.Round, "latency_ms": ms})
    return nil
}

// LoadLastState loads the last persisted state with pessimistic recovery
// (fallback to .bak if the main file is corrupt or missing).
func (fs *FileStore) LoadLastState(_ context.Context) (LastState, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    // Try main
    if s, err := readFile(fs.path); err == nil {
        metrics.Inc("state_recovery_total", map[string]string{"result": "ok"})
        logger.InfoJ("consensus_state", map[string]any{"op":"recovery", "result":"ok", "height": s.Height, "round": s.Round})
        return s, nil
    }
    // Try backup
    if s, err := readFile(fs.path + ".bak"); err == nil {
        metrics.Inc("state_recovery_total", map[string]string{"result": "fallback"})
        logger.InfoJ("consensus_state", map[string]any{"op":"recovery", "result":"fallback", "height": s.Height, "round": s.Round})
        return s, nil
    }
    metrics.Inc("state_recovery_total", map[string]string{"result": "fail"})
    logger.InfoJ("consensus_state", map[string]any{"op":"recovery", "result":"miss"})
    return LastState{}, ErrNotFound
}

// Close implements Store. For FileStore it is a no-op.
func (fs *FileStore) Close() error { return nil }

