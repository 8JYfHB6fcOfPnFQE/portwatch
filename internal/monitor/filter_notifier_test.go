package monitor_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func makeFilterEvent(port int, proto, addr string) alert.Event {
	return alert.NewEvent(ports.PortState{
		Port:  port,
		Proto: proto,
		Addr:  addr,
	}, "opened")
}

func TestFilterNotifier_ForwardsNonMatchingEvent(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error { got = ev; return nil })
	fn := monitor.NewFilterNotifier([]string{"blocked"}, next)

	ev := makeFilterEvent(8080, "tcp", "0.0.0.0")
	if err := fn.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Port != 8080 {
		t.Errorf("expected event forwarded, got port %d", got.Port)
	}
}

func TestFilterNotifier_DropsMatchingEvent(t *testing.T) {
	called := false
	next := alert.NotifierFunc(func(ev alert.Event) error { called = true; return nil })
	fn := monitor.NewFilterNotifier([]string{"tcp"}, next)

	ev := makeFilterEvent(9090, "tcp", "127.0.0.1")
	if err := fn.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected event to be dropped but next was called")
	}
}

func TestFilterNotifier_CaseInsensitiveMatch(t *testing.T) {
	called := false
	next := alert.NotifierFunc(func(ev alert.Event) error { called = true; return nil })
	fn := monitor.NewFilterNotifier([]string{"OPENED"}, next)

	ev := makeFilterEvent(443, "tcp", "0.0.0.0")
	if err := fn.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected case-insensitive match to drop event")
	}
}

func TestFilterNotifier_MatchesMetaValue(t *testing.T) {
	called := false
	next := alert.NotifierFunc(func(ev alert.Event) error { called = true; return nil })
	fn := monitor.NewFilterNotifier([]string{"internal"}, next)

	ev := makeFilterEvent(80, "tcp", "0.0.0.0")
	ev.Meta = map[string]string{"env": "internal"}
	if err := fn.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected meta-value match to drop event")
	}
}

func TestFilterNotifier_NilNext_NoError(t *testing.T) {
	fn := monitor.NewFilterNotifier([]string{"nomatch"}, nil)
	ev := makeFilterEvent(22, "tcp", "0.0.0.0")
	if err := fn.Send(ev); err != nil {
		t.Errorf("unexpected error with nil next: %v", err)
	}
}

func TestFilterNotifier_PropagatesNextError(t *testing.T) {
	want := errors.New("downstream failure")
	next := alert.NotifierFunc(func(ev alert.Event) error { return want })
	fn := monitor.NewFilterNotifier([]string{"nomatch"}, next)

	ev := makeFilterEvent(8080, "tcp", "0.0.0.0")
	if err := fn.Send(ev); !errors.Is(err, want) {
		t.Errorf("expected %v, got %v", want, err)
	}
}
