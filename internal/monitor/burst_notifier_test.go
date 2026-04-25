package monitor

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeBurstEvent(port int, proto string) alert.Event {
	return alert.NewEvent("opened", proto, port, "127.0.0.1")
}

func TestBurstNotifier_ForwardsWithinLimit(t *testing.T) {
	var received []alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error {
		received = append(received, ev)
		return nil
	})

	b := NewBurstNotifier(next, 3, time.Second)
	for i := 0; i < 3; i++ {
		_ = b.Send(makeBurstEvent(8080, "tcp"))
	}

	if len(received) != 3 {
		t.Fatalf("expected 3 forwarded events, got %d", len(received))
	}
}

func TestBurstNotifier_DropsAboveLimit(t *testing.T) {
	var received []alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error {
		received = append(received, ev)
		return nil
	})

	b := NewBurstNotifier(next, 2, time.Second)
	for i := 0; i < 5; i++ {
		_ = b.Send(makeBurstEvent(9090, "tcp"))
	}

	if len(received) != 2 {
		t.Fatalf("expected 2 forwarded events, got %d", len(received))
	}
}

func TestBurstNotifier_DifferentKeysIndependent(t *testing.T) {
	var received []alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error {
		received = append(received, ev)
		return nil
	})

	b := NewBurstNotifier(next, 1, time.Second)
	_ = b.Send(makeBurstEvent(80, "tcp"))
	_ = b.Send(makeBurstEvent(443, "tcp"))
	// second sends for each should be dropped
	_ = b.Send(makeBurstEvent(80, "tcp"))
	_ = b.Send(makeBurstEvent(443, "tcp"))

	if len(received) != 2 {
		t.Fatalf("expected 2 events (one per key), got %d", len(received))
	}
}

func TestBurstNotifier_ZeroLimit_DisablesBurst(t *testing.T) {
	var received []alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error {
		received = append(received, ev)
		return nil
	})

	b := NewBurstNotifier(next, 0, time.Second)
	for i := 0; i < 10; i++ {
		_ = b.Send(makeBurstEvent(22, "tcp"))
	}

	if len(received) != 10 {
		t.Fatalf("expected all 10 events forwarded when limit=0, got %d", len(received))
	}
}

func TestBurstNotifier_Reset_ClearsCounters(t *testing.T) {
	var received []alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error {
		received = append(received, ev)
		return nil
	})

	b := NewBurstNotifier(next, 1, time.Second)
	_ = b.Send(makeBurstEvent(8080, "tcp"))
	_ = b.Send(makeBurstEvent(8080, "tcp")) // dropped
	b.Reset()
	_ = b.Send(makeBurstEvent(8080, "tcp")) // allowed after reset

	if len(received) != 2 {
		t.Fatalf("expected 2 events after reset, got %d", len(received))
	}
}
