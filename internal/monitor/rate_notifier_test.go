package monitor

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeRateEvent(port int) alert.Event {
	e := alert.NewEvent("opened", port, "tcp")
	e.Message = "test"
	return e
}

func TestRateNotifier_ForwardsWithinLimit(t *testing.T) {
	var count int32
	next := alert.NotifierFunc(func(e alert.Event) error {
		atomic.AddInt32(&count, 1)
		return nil
	})
	rn := NewRateNotifier(3, time.Minute, next)
	for i := 0; i < 3; i++ {
		if err := rn.Send(makeRateEvent(8000 + i)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := atomic.LoadInt32(&count); got != 3 {
		t.Errorf("expected 3 forwarded, got %d", got)
	}
}

func TestRateNotifier_DropsAboveLimit(t *testing.T) {
	var count int32
	next := alert.NotifierFunc(func(e alert.Event) error {
		atomic.AddInt32(&count, 1)
		return nil
	})
	rn := NewRateNotifier(2, time.Minute, next)
	for i := 0; i < 5; i++ {
		_ = rn.Send(makeRateEvent(9000 + i))
	}
	if got := atomic.LoadInt32(&count); got != 2 {
		t.Errorf("expected 2 forwarded, got %d", got)
	}
	if rn.Dropped() != 3 {
		t.Errorf("expected 3 dropped, got %d", rn.Dropped())
	}
}

func TestRateNotifier_ResetsAfterWindow(t *testing.T) {
	var count int32
	next := alert.NotifierFunc(func(e alert.Event) error {
		atomic.AddInt32(&count, 1)
		return nil
	})
	rn := NewRateNotifier(1, 50*time.Millisecond, next)
	_ = rn.Send(makeRateEvent(7001))
	_ = rn.Send(makeRateEvent(7002)) // dropped
	time.Sleep(60 * time.Millisecond)
	_ = rn.Send(makeRateEvent(7003)) // new window, forwarded
	if got := atomic.LoadInt32(&count); got != 2 {
		t.Errorf("expected 2 forwarded across windows, got %d", got)
	}
}

func TestRateNotifier_ZeroLimit_DisablesRateLimit(t *testing.T) {
	var count int32
	next := alert.NotifierFunc(func(e alert.Event) error {
		atomic.AddInt32(&count, 1)
		return nil
	})
	rn := NewRateNotifier(0, time.Minute, next)
	for i := 0; i < 10; i++ {
		_ = rn.Send(makeRateEvent(6000 + i))
	}
	if got := atomic.LoadInt32(&count); got != 10 {
		t.Errorf("expected all 10 forwarded, got %d", got)
	}
}

func TestRateNotifier_NilNext_NoError(t *testing.T) {
	rn := NewRateNotifier(5, time.Minute, nil)
	if err := rn.Send(makeRateEvent(5000)); err != nil {
		t.Errorf("unexpected error with nil next: %v", err)
	}
}
