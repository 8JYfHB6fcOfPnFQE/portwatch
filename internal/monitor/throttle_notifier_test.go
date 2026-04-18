package monitor

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type captureNotifier struct {
	events []alert.Event
}

func (c *captureNotifier) Send(e alert.Event) error {
	c.events = append(c.events, e)
	return nil
}

func baseEvent(port int) alert.Event {
	return alert.NewEvent("opened", port, "tcp")
}

func TestThrottleNotifier_FirstAlertForwarded(t *testing.T) {
	cap := &captureNotifier{}
	tn := NewThrottleNotifier(cap, 5*time.Second)

	_ = tn.Send(baseEvent(8080))

	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
}

func TestThrottleNotifier_SecondAlertSuppressed(t *testing.T) {
	cap := &captureNotifier{}
	n := time.Now()
	tn := NewThrottleNotifier(cap, 5*time.Second)
	tn.nowFunc = func() time.Time { return n }

	_ = tn.Send(baseEvent(8080))
	_ = tn.Send(baseEvent(8080))

	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
}

func TestThrottleNotifier_AlertAfterWindowExpires(t *testing.T) {
	cap := &captureNotifier{}
	now := time.Now()
	tn := NewThrottleNotifier(cap, 5*time.Second)
	tn.nowFunc = func() time.Time { return now }

	_ = tn.Send(baseEvent(8080))

	now = now.Add(6 * time.Second)
	_ = tn.Send(baseEvent(8080))

	if len(cap.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(cap.events))
	}
}

func TestThrottleNotifier_DifferentPortsIndependent(t *testing.T) {
	cap := &captureNotifier{}
	tn := NewThrottleNotifier(cap, 5*time.Second)

	_ = tn.Send(baseEvent(8080))
	_ = tn.Send(baseEvent(9090))

	if len(cap.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(cap.events))
	}
}

func TestThrottleNotifier_Reset_AllowsImmediateResend(t *testing.T) {
	cap := &captureNotifier{}
	now := time.Now()
	tn := NewThrottleNotifier(cap, 5*time.Second)
	tn.nowFunc = func() time.Time { return now }

	_ = tn.Send(baseEvent(8080))
	tn.Reset("tcp", 8080)
	_ = tn.Send(baseEvent(8080))

	if len(cap.events) != 2 {
		t.Fatalf("expected 2 events after reset, got %d", len(cap.events))
	}
}
