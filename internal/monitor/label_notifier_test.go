package monitor_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func makeBaseEvent() alert.Event {
	return alert.NewEvent("opened", "tcp", 8080, "127.0.0.1")
}

func TestLabelNotifier_AttachesLabels(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(_ context.Context, ev alert.Event) error {
		got = ev
		return nil
	})

	ln := monitor.NewLabelNotifier(map[string]string{"env": "prod", "region": "us-east"}, next)
	ev := makeBaseEvent()
	if err := ln.Send(context.Background(), ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Meta["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", got.Meta["env"])
	}
	if got.Meta["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %q", got.Meta["region"])
	}
}

func TestLabelNotifier_OverridesExistingMeta(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(_ context.Context, ev alert.Event) error {
		got = ev
		return nil
	})

	ln := monitor.NewLabelNotifier(map[string]string{"env": "staging"}, next)
	ev := makeBaseEvent()
	ev.Meta = map[string]string{"env": "dev", "host": "box1"}

	if err := ln.Send(context.Background(), ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Meta["env"] != "staging" {
		t.Errorf("label should override existing meta: got %q", got.Meta["env"])
	}
	if got.Meta["host"] != "box1" {
		t.Errorf("unrelated meta key should be preserved: got %q", got.Meta["host"])
	}
}

func TestLabelNotifier_EmptyLabels_PassesThrough(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(_ context.Context, ev alert.Event) error {
		got = ev
		return nil
	})

	ln := monitor.NewLabelNotifier(nil, next)
	ev := makeBaseEvent()
	ev.Meta = map[string]string{"k": "v"}

	if err := ln.Send(context.Background(), ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Meta["k"] != "v" {
		t.Errorf("expected original meta preserved, got %v", got.Meta)
	}
}

func TestLabelNotifier_NilNext_NoError(t *testing.T) {
	ln := monitor.NewLabelNotifier(map[string]string{"x": "y"}, nil)
	if err := ln.Send(context.Background(), makeBaseEvent()); err != nil {
		t.Fatalf("nil next should not cause error: %v", err)
	}
}

func TestLabelNotifier_DoesNotMutateOriginalEvent(t *testing.T) {
	next := alert.NotifierFunc(func(_ context.Context, _ alert.Event) error { return nil })
	ln := monitor.NewLabelNotifier(map[string]string{"env": "prod"}, next)

	ev := makeBaseEvent()
	ev.Meta = map[string]string{"original": "true"}
	origLen := len(ev.Meta)

	_ = ln.Send(context.Background(), ev)

	if len(ev.Meta) != origLen {
		t.Errorf("original event meta was mutated")
	}
}
