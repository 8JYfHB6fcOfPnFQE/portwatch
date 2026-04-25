package monitor

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makePriorityEvent(msg string) alert.Event {
	return alert.NewEvent("opened", 8080, "tcp", msg)
}

func TestPriorityNotifier_DefaultPriority(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })
	n := NewPriorityNotifier(next, map[string][]string{}, "medium")

	_ = n.Send(makePriorityEvent("some random message"))

	if got.Meta["priority"] != "medium" {
		t.Fatalf("expected medium, got %q", got.Meta["priority"])
	}
}

func TestPriorityNotifier_MatchesHighKeyword(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })
	rules := map[string][]string{
		"high": {"ssh", "rdp"},
	}
	n := NewPriorityNotifier(next, rules, "low")

	_ = n.Send(makePriorityEvent("port opened: ssh service detected"))

	if got.Meta["priority"] != "high" {
		t.Fatalf("expected high, got %q", got.Meta["priority"])
	}
}

func TestPriorityNotifier_CriticalBeatsHigh(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })
	rules := map[string][]string{
		"high":     {"ssh"},
		"critical": {"root"},
	}
	n := NewPriorityNotifier(next, rules, "low")

	// Message contains both; critical should win due to ordering.
	_ = n.Send(makePriorityEvent("root ssh access detected"))

	if got.Meta["priority"] != "critical" {
		t.Fatalf("expected critical, got %q", got.Meta["priority"])
	}
}

func TestPriorityNotifier_CaseInsensitiveMatch(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })
	rules := map[string][]string{
		"high": {"SSH"},
	}
	n := NewPriorityNotifier(next, rules, "low")

	_ = n.Send(makePriorityEvent("unexpected ssh port opened"))

	if got.Meta["priority"] != "high" {
		t.Fatalf("expected high, got %q", got.Meta["priority"])
	}
}

func TestPriorityNotifier_NilNext_NoError(t *testing.T) {
	n := NewPriorityNotifier(nil, map[string][]string{}, "low")
	if err := n.Send(makePriorityEvent("test")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPriorityNotifier_EmptyDefaultFallsToLow(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })
	n := NewPriorityNotifier(next, map[string][]string{}, "")

	_ = n.Send(makePriorityEvent("nothing special"))

	if got.Meta["priority"] != "low" {
		t.Fatalf("expected low, got %q", got.Meta["priority"])
	}
}
