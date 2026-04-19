package monitor

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type countingNotifier struct {
	calls  int
	failN  int // fail the first N calls
	last   alert.Event
}

func (c *countingNotifier) Send(e alert.Event) error {
	c.calls++
	c.last = e
	if c.calls <= c.failN {
		return errors.New("transient error")
	}
	return nil
}

func TestRetryNotifier_SucceedsFirstTry(t *testing.T) {
	inner := &countingNotifier{}
	r := NewRetryNotifier(inner, 3, time.Millisecond, nil)
	if err := r.Send(alert.NewEvent("opened", 80, "tcp")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls)
	}
}

func TestRetryNotifier_RetriesOnFailure(t *testing.T) {
	inner := &countingNotifier{failN: 2}
	r := NewRetryNotifier(inner, 3, time.Millisecond, nil)
	if err := r.Send(alert.NewEvent("opened", 443, "tcp")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryNotifier_ReturnsLastError(t *testing.T) {
	inner := &countingNotifier{failN: 5}
	r := NewRetryNotifier(inner, 3, time.Millisecond, nil)
	if err := r.Send(alert.NewEvent("opened", 22, "tcp")); err == nil {
		t.Fatal("expected error but got nil")
	}
	if inner.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryNotifier_LogsAttempts(t *testing.T) {
	inner := &countingNotifier{failN: 1}
	logger := log.New(os.Stderr, "", 0)
	r := NewRetryNotifier(inner, 2, time.Millisecond, logger)
	if err := r.Send(alert.NewEvent("opened", 8080, "tcp")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetryNotifier_MinOneTry(t *testing.T) {
	inner := &countingNotifier{}
	r := NewRetryNotifier(inner, 0, time.Millisecond, nil)
	_ = r.Send(alert.NewEvent("opened", 9090, "tcp"))
	if inner.calls != 1 {
		t.Fatalf("expected at least 1 call, got %d", inner.calls)
	}
}
