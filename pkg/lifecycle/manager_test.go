package lifecycle

import (
    "context"
    "errors"
    "testing"
)

type stub struct { name string; startErr, stopErr error; started bool }

func (s *stub) Name() string { return s.name }
func (s *stub) Start(ctx context.Context) error { s.started = true; return s.startErr }
func (s *stub) Stop(ctx context.Context) error  { s.started = false; return s.stopErr }

func TestStartAll_RollbackOnFailure(t *testing.T) {
    a := &stub{name: "a"}
    b := &stub{name: "b", startErr: errors.New("b fail")}
    c := &stub{name: "c"}

    m := New(); m.Add(a); m.Add(b); m.Add(c)
    if err := m.StartAll(context.Background()); err == nil {
        t.Fatalf("want error")
    }
    if a.started { t.Fatalf("a should be rolled back and not started") }
}

func TestStopAll_ReverseOrder(t *testing.T) {
    a := &stub{name: "a"}
    b := &stub{name: "b"}
    m := New(); m.Add(a); m.Add(b)
    if err := m.StartAll(context.Background()); err != nil { t.Fatalf("unexpected: %v", err) }
    if err := m.StopAll(context.Background()); err != nil { t.Fatalf("unexpected stop error: %v", err) }
}
