package monitor

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeOpsEvent(port int) alert.Event {
	return alert.NewEvent("opened", port, "tcp", "")
}

func TestOpsNotifier_CountsSentEvents(t *testing.T) {
	var received int
	next := alert.NotifierFunc(func(_ context.Context, _ alert.Event) error {
		received++
		return nil
	})
	m := &OpsMetrics{}
	on := NewOpsNotifier(next, m, log.New(os.Stderr, "", 0))

	_ = on.Send(context.Background(), makeOpsEvent(8080))
	_ = on.Send(context.Background(), makeOpsEvent(9090))

	if m.Sent != 2 {
		t.Fatalf("expected 2 sent, got %d", m.Sent)
	}
	if m.Failed != 0 {
		t.Fatalf("expected 0 failed, got %d", m.Failed)
	}
	if received != 2 {
		t.Fatalf("expected next called twice, got %d", received)
	}
}

func TestOpsNotifier_CountsFailedEvents(t *testing.T) {
	next := alert.NotifierFunc(func(_ context.Context, _ alert.Event) error {
		return errors.New("downstream error")
	})
	m := &OpsMetrics{}
	on := NewOpsNotifier(next, m, nil)

	err := on.Send(context.Background(), makeOpsEvent(443))
	if err == nil {
		t.Fatal("expected error")
	}
	if m.Failed != 1 {
		t.Fatalf("expected 1 failed, got %d", m.Failed)
	}
	if m.Sent != 0 {
		t.Fatalf("expected 0 sent, got %d", m.Sent)
	}
}

func TestOpsNotifier_AvgLatency_NoEvents(t *testing.T) {
	next := alert.NotifierFunc(func(_ context.Context, _ alert.Event) error { return nil })
	on := NewOpsNotifier(next, nil, nil)
	if on.AvgLatencyMs() != 0 {
		t.Fatal("expected 0 avg latency with no events")
	}
}

func TestOpsNotifier_AvgLatency_AfterSend(t *testing.T) {
	next := alert.NotifierFunc(func(_ context.Context, _ alert.Event) error { return nil })
	on := NewOpsNotifier(next, nil, nil)
	_ = on.Send(context.Background(), makeOpsEvent(80))
	if on.AvgLatencyMs() < 0 {
		t.Fatal("avg latency must be non-negative")
	}
}

func TestOpsNotifier_NilMetrics_Initialised(t *testing.T) {
	next := alert.NotifierFunc(func(_ context.Context, _ alert.Event) error { return nil })
	on := NewOpsNotifier(next, nil, nil)
	if on.metrics == nil {
		t.Fatal("metrics must not be nil after construction")
	}
}
