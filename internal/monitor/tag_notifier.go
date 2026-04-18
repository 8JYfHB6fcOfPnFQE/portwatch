package monitor

import (
	"context"

	"github.com/user/portwatch/internal/alert"
)

// TagNotifier wraps another Notifier and drops events whose tags match the
// configured suppression list.
type TagNotifier struct {
	filter *TagFilter
	next   alert.Notifier
}

// NewTagNotifier returns a TagNotifier that suppresses events tagged with any
// of suppressedTags before forwarding to next.
func NewTagNotifier(suppressedTags []string, next alert.Notifier) *TagNotifier {
	return &TagNotifier{
		filter: NewTagFilter(suppressedTags),
		next:   next,
	}
}

// Send drops the event when its tags are suppressed; otherwise it delegates.
func (n *TagNotifier) Send(ctx context.Context, ev alert.Event) error {
	if n.filter.IsSuppressed(ev.Tags) {
		return nil
	}
	return n.next.Send(ctx, ev)
}
