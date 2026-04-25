package monitor

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

type captureNotifier struct {
	events []alert.Event
	err    error
}

func (c *captureNotifier) Send(e alert.Event) error {
	c.events = append(c.events, e)
	return c.err
}

func makeRoutingEvent(proto, action string, meta map[string]string) alert.Event {
	e := alert.NewEvent(8080, proto, action)
	e.Meta = meta
	return e
}

func TestRoutingNotifier_MatchesProto(t *testing.T) {
	tcp := &captureNotifier{}
	udp := &captureNotifier{}
	rn := NewRoutingNotifier([]RouteRule{
		{Field: "proto", Value: "tcp", Notifier: tcp},
		{Field: "proto", Value: "udp", Notifier: udp},
	}, nil)

	_ = rn.Send(makeRoutingEvent("tcp", "opened", nil))
	_ = rn.Send(makeRoutingEvent("udp", "opened", nil))

	if len(tcp.events) != 1 {
		t.Fatalf("expected 1 tcp event, got %d", len(tcp.events))
	}
	if len(udp.events) != 1 {
		t.Fatalf("expected 1 udp event, got %d", len(udp.events))
	}
}

func TestRoutingNotifier_MatchesAction(t *testing.T) {
	opened := &captureNotifier{}
	rn := NewRoutingNotifier([]RouteRule{
		{Field: "action", Value: "opened", Notifier: opened},
	}, nil)

	_ = rn.Send(makeRoutingEvent("tcp", "opened", nil))
	_ = rn.Send(makeRoutingEvent("tcp", "closed", nil))

	if len(opened.events) != 1 {
		t.Fatalf("expected 1 opened event, got %d", len(opened.events))
	}
}

func TestRoutingNotifier_FallbackOnNoMatch(t *testing.T) {
	fallback := &captureNotifier{}
	rn := NewRoutingNotifier([]RouteRule{
		{Field: "proto", Value: "udp", Notifier: &captureNotifier{}},
	}, fallback)

	_ = rn.Send(makeRoutingEvent("tcp", "opened", nil))

	if len(fallback.events) != 1 {
		t.Fatalf("expected fallback to receive event, got %d", len(fallback.events))
	}
}

func TestRoutingNotifier_NoFallback_NoError(t *testing.T) {
	rn := NewRoutingNotifier([]RouteRule{}, nil)
	if err := rn.Send(makeRoutingEvent("tcp", "opened", nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoutingNotifier_MatchesMetaKey(t *testing.T) {
	dest := &captureNotifier{}
	rn := NewRoutingNotifier([]RouteRule{
		{Field: "env", Value: "prod", Notifier: dest},
	}, nil)

	_ = rn.Send(makeRoutingEvent("tcp", "opened", map[string]string{"env": "production"}))
	_ = rn.Send(makeRoutingEvent("tcp", "opened", map[string]string{"env": "staging"}))

	if len(dest.events) != 1 {
		t.Fatalf("expected 1 meta-matched event, got %d", len(dest.events))
	}
}

func TestRoutingNotifier_PropagatesError(t *testing.T) {
	fail := &captureNotifier{err: errors.New("send failed")}
	rn := NewRoutingNotifier([]RouteRule{
		{Field: "proto", Value: "tcp", Notifier: fail},
	}, nil)

	err := rn.Send(makeRoutingEvent("tcp", "opened", nil))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
