package monitor

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeTransformEvent() alert.Event {
	return alert.Event{
		Port:   8080,
		Proto:  "tcp",
		Action: "opened",
		Addr:   "192.168.1.5",
		At:     time.Now(),
		Meta:   map[string]string{"env": "prod"},
	}
}

func TestTransformNotifier_UpperCaseAction(t *testing.T) {
	var got alert.Event
	n := NewTransformNotifier(
		alert.NotifierFunc(func(e alert.Event) error { got = e; return nil }),
		UpperCaseAction(),
	)
	_ = n.Send(makeTransformEvent())
	if got.Action != "OPENED" {
		t.Errorf("expected OPENED, got %q", got.Action)
	}
}

func TestTransformNotifier_SetMeta(t *testing.T) {
	var got alert.Event
	n := NewTransformNotifier(
		alert.NotifierFunc(func(e alert.Event) error { got = e; return nil }),
		SetMeta("region", "us-east-1"),
	)
	_ = n.Send(makeTransformEvent())
	if got.Meta["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", got.Meta["region"])
	}
	// original key preserved
	if got.Meta["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", got.Meta["env"])
	}
}

func TestTransformNotifier_RedactAddr(t *testing.T) {
	var got alert.Event
	n := NewTransformNotifier(
		alert.NotifierFunc(func(e alert.Event) error { got = e; return nil }),
		RedactAddr(),
	)
	_ = n.Send(makeTransformEvent())
	if got.Addr != "" {
		t.Errorf("expected empty addr, got %q", got.Addr)
	}
}

func TestTransformNotifier_ChainedTransforms(t *testing.T) {
	var got alert.Event
	n := NewTransformNotifier(
		alert.NotifierFunc(func(e alert.Event) error { got = e; return nil }),
		UpperCaseAction(),
		RedactAddr(),
		SetMeta("transformed", "yes"),
	)
	_ = n.Send(makeTransformEvent())
	if got.Action != "OPENED" {
		t.Errorf("action: want OPENED, got %q", got.Action)
	}
	if got.Addr != "" {
		t.Errorf("addr: want empty, got %q", got.Addr)
	}
	if got.Meta["transformed"] != "yes" {
		t.Errorf("meta: want yes, got %q", got.Meta["transformed"])
	}
}

func TestTransformNotifier_NilNext_NoError(t *testing.T) {
	n := NewTransformNotifier(nil, UpperCaseAction())
	if err := n.Send(makeTransformEvent()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTransformNotifier_PropagatesNextError(t *testing.T) {
	sentinel := errors.New("downstream failure")
	n := NewTransformNotifier(
		alert.NotifierFunc(func(_ alert.Event) error { return sentinel }),
	)
	if err := n.Send(makeTransformEvent()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}
