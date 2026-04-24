package monitor

import (
	"context"

	"github.com/user/portwatch/internal/alert"
)

// LabelNotifier attaches static key/value labels to every event before
// forwarding it downstream. Labels are merged into the event's metadata map,
// with notifier-supplied values taking precedence over existing ones.
type LabelNotifier struct {
	labels map[string]string
	next   alert.Notifier
}

// NewLabelNotifier returns a LabelNotifier that stamps each event with the
// provided labels and forwards the enriched event to next.
func NewLabelNotifier(labels map[string]string, next alert.Notifier) *LabelNotifier {
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return &LabelNotifier{labels: copy, next: next}
}

// Send attaches the configured labels to ev and forwards the result.
func (l *LabelNotifier) Send(ctx context.Context, ev alert.Event) error {
	if len(l.labels) > 0 {
		merged := make(map[string]string, len(ev.Meta)+len(l.labels))
		for k, v := range ev.Meta {
			merged[k] = v
		}
		for k, v := range l.labels {
			merged[k] = v
		}
		ev.Meta = merged
	}
	if l.next == nil {
		return nil
	}
	return l.next.Send(ctx, ev)
}
