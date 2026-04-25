package monitor_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

type countingNotifier struct {
	calls atomic.Int64
	err   error
}

func (c *countingNotifier) Send(_ alert.Event) error {
	c.calls.Add(1)
	return c.err
}

func makeSamplingEvent() alert.Event {
	return alert.NewEvent("opened", "tcp", "127.0.0.1", 9090)
}

func TestSamplingNotifier_RateOne_ForwardsAll(t *testing.T) {
	next := &countingNotifier{}
	sn := monitor.NewSamplingNotifier(1.0, next)

	for i := 0; i < 100; i++ {
		_ = sn.Send(makeSamplingEvent())
	}
	if next.calls.Load() != 100 {
		t.Fatalf("expected 100 forwarded, got %d", next.calls.Load())
	}
}

func TestSamplingNotifier_RateZero_DropsAll(t *testing.T) {
	next := &countingNotifier{}
	sn := monitor.NewSamplingNotifier(0.0, next)

	for i := 0; i < 100; i++ {
		_ = sn.Send(makeSamplingEvent())
	}
	if next.calls.Load() != 0 {
		t.Fatalf("expected 0 forwarded, got %d", next.calls.Load())
	}
}

func TestSamplingNotifier_ClampsBelowZero(t *testing.T) {
	sn := monitor.NewSamplingNotifier(-5.0, nil)
	if sn.Rate() != 0 {
		t.Fatalf("expected rate clamped to 0, got %g", sn.Rate())
	}
}

func TestSamplingNotifier_ClampsAboveOne(t *testing.T) {
	sn := monitor.NewSamplingNotifier(3.14, nil)
	if sn.Rate() != 1.0 {
		t.Fatalf("expected rate clamped to 1.0, got %g", sn.Rate())
	}
}

func TestSamplingNotifier_NilNext_NoError(t *testing.T) {
	sn := monitor.NewSamplingNotifier(1.0, nil)
	if err := sn.Send(makeSamplingEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSamplingNotifier_PropagatesError(t *testing.T) {
	want := errors.New("downstream failure")
	next := &countingNotifier{err: want}
	sn := monitor.NewSamplingNotifier(1.0, next)

	if err := sn.Send(makeSamplingEvent()); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestSamplingNotifier_PartialRate_SamplesRoughly(t *testing.T) {
	next := &countingNotifier{}
	sn := monitor.NewSamplingNotifier(0.5, next)

	const n = 10_000
	for i := 0; i < n; i++ {
		_ = sn.Send(makeSamplingEvent())
	}
	got := next.calls.Load()
	// Allow generous tolerance: expect between 30% and 70% forwarded.
	if got < 3000 || got > 7000 {
		t.Fatalf("rate 0.5 over %d events: got %d forwarded (expected ~5000)", n, got)
	}
}
