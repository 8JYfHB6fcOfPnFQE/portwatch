package monitor

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type captureNotifier struct {
	mu     sync.Mutex
	events []alert.Event
}

func (c *captureNotifier) Send(e alert.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, e)
	return nil
}

func (c *captureNotifier) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.events)
}

func TestDedupNotifier_DropsSecondIdenticalEvent(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDedupNotifier(cap, 5*time.Second)
	e := alert.NewEvent("opened", 8080, "tcp")

	_ = d.Send(e)
	_ = d.Send(e)

	if cap.count() != 1 {
		t.Fatalf("expected 1 event, got %d", cap.count())
	}
}

func TestDedupNotifier_AllowsAfterTTLExpires(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDedupNotifier(cap, 10*time.Millisecond)
	e := alert.NewEvent("opened", 8080, "tcp")

	_ = d.Send(e)
	time.Sleep(20 * time.Millisecond)
	_ = d.Send(e)

	if cap.count() != 2 {
		t.Fatalf("expected 2 events, got %d", cap.count())
	}
}

func TestDedupNotifier_DifferentPortsIndependent(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDedupNotifier(cap, 5*time.Second)

	_ = d.Send(alert.NewEvent("opened", 8080, "tcp"))
	_ = d.Send(alert.NewEvent("opened", 9090, "tcp"))

	if cap.count() != 2 {
		t.Fatalf("expected 2 events, got %d", cap.count())
	}
}

func TestDedupNotifier_Flush_ResetsState(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDedupNotifier(cap, 5*time.Second)
	e := alert.NewEvent("opened", 8080, "tcp")

	_ = d.Send(e)
	d.Flush()
	_ = d.Send(e)

	if cap.count() != 2 {
		t.Fatalf("expected 2 events after flush, got %d", cap.count())
	}
}
