package lifecycle

import (
    "context"
    "errors"
    "testing"
)

type mockSvc struct {
    name     string
    startErr error
    stopErr  error
    calls    *[]string
}

func (m *mockSvc) Name() string { return m.name }
func (m *mockSvc) Start(ctx context.Context) error {
    if m.calls != nil { *m.calls = append(*m.calls, "start:"+m.name) }
    return m.startErr
}
func (m *mockSvc) Stop(ctx context.Context) error {
    if m.calls != nil { *m.calls = append(*m.calls, "stop:"+m.name) }
    return m.stopErr
}

func TestStartAllSuccessAndStopReverse(t *testing.T) {
    var calls []string
    s1 := &mockSvc{name: "s1", calls: &calls}
    s2 := &mockSvc{name: "s2", calls: &calls}

    m := New()
    m.Add(s1)
    m.Add(s2)

    if err := m.StartAll(context.Background()); err != nil {
        t.Fatalf("unexpected start error: %v", err)
    }

    if err := m.StopAll(context.Background()); err != nil {
        t.Fatalf("unexpected stop error: %v", err)
    }

    want := []string{"start:s1", "start:s2", "stop:s2", "stop:s1"}
    if len(calls) != len(want) {
        t.Fatalf("calls len=%d want=%d: %+v", len(calls), len(want), calls)
    }
    for i := range want {
        if calls[i] != want[i] {
            t.Fatalf("calls[%d]=%s want %s", i, calls[i], want[i])
        }
    }
}

func TestStartAllRollbackOnFailure(t *testing.T) {
    var calls []string
    s1 := &mockSvc{name: "s1", calls: &calls}
    boom := errors.New("start-fail")
    s2 := &mockSvc{name: "s2", calls: &calls, startErr: boom}

    m := New()
    m.Add(s1)
    m.Add(s2)

    err := m.StartAll(context.Background())
    if err == nil {
        t.Fatalf("want error")
    }
    if !errors.Is(err, boom) {
        t.Fatalf("want err to contain start-fail; got %v", err)
    }

    want := []string{"start:s1", "start:s2", "stop:s1"}
    if len(calls) != len(want) {
        t.Fatalf("calls len=%d want=%d: %+v", len(calls), len(want), calls)
    }
    for i := range want {
        if calls[i] != want[i] {
            t.Fatalf("calls[%d]=%s want %s", i, calls[i], want[i])
        }
    }
}

func TestStopAllAggregatesErrorsAndReverseOrder(t *testing.T) {
    var calls []string
    e1 := errors.New("stop1")
    e2 := errors.New("stop2")
    s1 := &mockSvc{name: "s1", calls: &calls, stopErr: e1}
    s2 := &mockSvc{name: "s2", calls: &calls, stopErr: e2}

    m := New()
    m.Add(s1)
    m.Add(s2)

    if err := m.StartAll(context.Background()); err != nil {
        t.Fatalf("unexpected start error: %v", err)
    }

    err := m.StopAll(context.Background())
    if err == nil {
        t.Fatalf("want aggregated error")
    }
    if !errors.Is(err, e1) || !errors.Is(err, e2) {
        t.Fatalf("aggregated error should include stop1 and stop2; got %v", err)
    }

    wantOrder := []string{"start:s1", "start:s2", "stop:s2", "stop:s1"}
    if len(calls) != len(wantOrder) {
        t.Fatalf("calls len=%d want=%d: %+v", len(calls), len(wantOrder), calls)
    }
    for i := range wantOrder {
        if calls[i] != wantOrder[i] {
            t.Fatalf("calls[%d]=%s want %s", i, calls[i], wantOrder[i])
        }
    }
}