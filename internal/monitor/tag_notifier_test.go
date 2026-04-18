package monitor

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

type captureNotifier struct {
	events []alert.Event
}

func (c *captureNotifier) Send(_ context.Context, ev alert.Event) error {
	c.events = append(c.events, ev)
	return nil
}

func TestTagNotifier_DropsTaggedEvent(t *testing.T) {
	cap := &captureNotifier{}
	n := NewTagNotifier([]string{"infra"}, cap)

	ev := alert.Event{Tags: []string{"infra"}}
	_ = n.Send(context.Background(), ev)

	if len(cap.events) != 0 {
		t.Errorf("expected 0 forwarded events, got %d", len(cap.events))
	}
}

func TestTagNotifier_ForwardsUntaggedEvent(t *testing.T) {
	cap := &captureNotifier{}
	n := NewTagNotifier([]string{"infra"}, cap)

	ev := alert.Event{Tags: []string{"prod"}}
	_ = n.Send(context.Background(), ev)

	if len(cap.events) != 1 {
		t.Errorf("expected 1 forwarded event, got %d", len(cap.events))
	}
}

func TestTagNotifier_ForwardsWhenNoSuppressedTags(t *testing.T) {
	cap := &captureNotifier{}
	n := NewTagNotifier([]string{}, cap)

	ev := alert.Event{Tags: []string{"infra", "dev"}}
	_ = n.Send(context.Background(), ev)

	if len(cap.events) != 1 {
		t.Errorf("expected 1 forwarded event, got %d", len(cap.events))
	}
}
