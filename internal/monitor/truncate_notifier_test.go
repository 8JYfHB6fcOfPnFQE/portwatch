package monitor

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeTruncateEvent(msg, addr string) alert.Event {
	e := alert.NewEvent("opened", 8080, "tcp")
	e.Message = msg
	e.Addr = addr
	return e
}

func TestTruncateNotifier_ShortMessage_Unchanged(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })

	n := NewTruncateNotifier(50, next)
	e := makeTruncateEvent("short", "127.0.0.1")
	if err := n.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Message != "short" {
		t.Errorf("expected 'short', got %q", got.Message)
	}
}

func TestTruncateNotifier_LongMessage_Truncated(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })

	n := NewTruncateNotifier(10, next)
	e := makeTruncateEvent("this message is definitely too long", "127.0.0.1")
	if err := n.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len([]rune(got.Message)) > 10 {
		t.Errorf("message not truncated: %q (len %d)", got.Message, len(got.Message))
	}
	if got.Message[len(got.Message)-3:] != "..." {
		t.Errorf("expected ellipsis suffix, got %q", got.Message)
	}
}

func TestTruncateNotifier_TruncatesAddr(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })

	n := NewTruncateNotifier(8, next)
	e := makeTruncateEvent("ok", "192.168.100.200")
	_ = n.Send(e)
	if len([]rune(got.Addr)) > 8 {
		t.Errorf("addr not truncated: %q", got.Addr)
	}
}

func TestTruncateNotifier_TruncatesMetaValues(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })

	n := NewTruncateNotifier(6, next)
	e := makeTruncateEvent("hi", "127.0.0.1")
	e.Meta = map[string]string{"note": "a very long meta value"}
	_ = n.Send(e)
	if len([]rune(got.Meta["note"])) > 6 {
		t.Errorf("meta value not truncated: %q", got.Meta["note"])
	}
}

func TestTruncateNotifier_ZeroMaxLen_ForwardsUnchanged(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })

	n := NewTruncateNotifier(0, next)
	const long = "this should not be truncated at all because maxLen is zero"
	e := makeTruncateEvent(long, "127.0.0.1")
	_ = n.Send(e)
	if got.Message != long {
		t.Errorf("expected unchanged message, got %q", got.Message)
	}
}

func TestTruncateNotifier_NilNext_NoError(t *testing.T) {
	n := NewTruncateNotifier(20, nil)
	e := makeTruncateEvent("hello", "127.0.0.1")
	if err := n.Send(e); err != nil {
		t.Errorf("expected nil error with nil next, got %v", err)
	}
}
