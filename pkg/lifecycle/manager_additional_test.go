package lifecycle

import (
    "context"
    "errors"
    "testing"
)

// Empty manager should StartAll/StopAll without error.
func TestManager_Empty_StartStopOK(t *testing.T) {
    m := New()
    if err := m.StartAll(context.Background()); err != nil {
        t.Fatalf("unexpected start error: %v", err)
    }
    if err := m.StopAll(context.Background()); err != nil {
        t.Fatalf("unexpected stop error: %v", err)
    }
}

// When second service fails to start, ensure third is never started and first is rolled back.
func TestStartAll_FailureShortCircuitsAndRollback(t *testing.T) {
    var calls []string
    s1 := &mockSvc{name: "s1", calls: &calls}
    boom := errors.New("start-fail-s2")
    s2 := &mockSvc{name: "s2", calls: &calls, startErr: boom}
    s3 := &mockSvc{name: "s3", calls: &calls}

    m := New()
    m.Add(s1)
    m.Add(s2)
    m.Add(s3)

    err := m.StartAll(context.Background())
    if err == nil || !errors.Is(err, boom) {
        t.Fatalf("expected start error including boom, got %v", err)
    }
    // Expect s3 was never started; s1 started then rolled back (stopped)
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

// Rollback errors should be aggregated with the start error via errors.Join.
func TestStartAll_RollbackAggregatesErrors(t *testing.T) {
    var calls []string
    stop1 := errors.New("stop1")
    s1 := &mockSvc{name: "s1", calls: &calls, stopErr: stop1}
    boom := errors.New("start-fail-s2")
    s2 := &mockSvc{name: "s2", calls: &calls, startErr: boom}

    m := New()
    m.Add(s1)
    m.Add(s2)

    err := m.StartAll(context.Background())
    if err == nil {
        t.Fatalf("want aggregated error")
    }
    if !errors.Is(err, boom) || !errors.Is(err, stop1) {
        t.Fatalf("aggregated error should include start and rollback stop errors; got %v", err)
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