package state

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestFileStore_SaveLoad_OK(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "laststate.dat")
    fs := NewFileStore(path)
    want := LastState{Height: 42, Round: 7}
    if err := fs.SaveLastState(context.Background(), want); err != nil {
        t.Fatalf("save: %v", err)
    }
    got, err := fs.LoadLastState(context.Background())
    if err != nil { t.Fatalf("load: %v", err) }
    if got != want { t.Fatalf("state mismatch: got=%+v want=%+v", got, want) }
}

func TestFileStore_Load_FallbackToBackup_OnCorruption(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "laststate.dat")
    fs := NewFileStore(path)
    // Save v1 (no backup yet)
    if err := fs.SaveLastState(context.Background(), LastState{Height: 1, Round: 1}); err != nil {
        t.Fatalf("save1: %v", err)
    }
    // Save v2 (creates .bak of v1)
    if err := fs.SaveLastState(context.Background(), LastState{Height: 2, Round: 2}); err != nil {
        t.Fatalf("save2: %v", err)
    }
    // Corrupt main file by truncation
    if err := os.Truncate(path, 8); err != nil { t.Fatalf("truncate: %v", err) }
    got, err := fs.LoadLastState(context.Background())
    if err != nil { t.Fatalf("load after corrupt: %v", err) }
    if got.Height != 1 || got.Round != 1 {
        t.Fatalf("fallback mismatch: got=%+v want={1 1}", got)
    }
}

func TestFileStore_Load_NotFound(t *testing.T) {
    dir := t.TempDir()
    fs := NewFileStore(filepath.Join(dir, "missing.dat"))
    if _, err := fs.LoadLastState(context.Background()); err == nil {
        t.Fatalf("want ErrNotFound")
    }
}

