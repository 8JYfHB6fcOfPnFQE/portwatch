package monitor_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

type digestSink struct {
	mu     sync.Mutex
	events []alert.Event
}

func (s *digestSink) Send(e alert.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
	return nil
}

func (s *digestSink) last() alert.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.events) == 0 {
		return alert.Event{}
	}
	return s.events[len(s.events)-1]
}

func TestDigestNotifier_BatchesEvents(t *testing.T) {
	sink := &digestSink{}
	d := monitor.NewDigestNotifier(80*time.Millisecond, sink)
	defer d.Stop()

	_ = d.Send(alert.NewEvent("opened", 8080, "tcp"))
	_ = d.Send(alert.NewEvent("opened", 9090, "tcp"))

	time.Sleep(150 * time.Millisecond)

	e := sink.last()
	if e.Kind != "digest" {
		t.Fatalf("expected digest event, got %q", e.Kind)
	}
	if !strings.Contains(e.Message, "2 event(s)") {
		t.Errorf("expected 2 events in digest, got: %s", e.Message)
	}
}

func TestDigestNotifier_EmptyWindow_NoFlush(t *testing.T) {
	sink := &digestSink{}
	d := monitor.NewDigestNotifier(50*time.Millisecond, sink)
	defer d.Stop()

	time.Sleep(120 * time.Millisecond)

	sink.mu.Lock()
	n := len(sink.events)
	sink.mu.Unlock()
	if n != 0 {
		t.Errorf("expected no events for empty window, got %d", n)
	}
}

func TestDigestNotifier_Stop_FlushesRemaining(t *testing.T) {
	sink := &digestSink{}
	d := monitor.NewDigestNotifier(10*time.Second, sink)

	_ = d.Send(alert.NewEvent("closed", 443, "tcp"))
	d.Stop()

	e := sink.last()
	if e.Kind != "digest" {
		t.Fatalf("expected digest on stop, got %q", e.Kind)
	}
	if !strings.Contains(e.Message, "port=443") {
		t.Errorf("expected port 443 in digest message: %s", e.Message)
	}
}
