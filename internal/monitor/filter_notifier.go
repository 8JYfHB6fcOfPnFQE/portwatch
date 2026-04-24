package monitor

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// FilterNotifier drops events whose message or metadata fields match any of
// the configured substring patterns before forwarding to next.
type FilterNotifier struct {
	patterns []string
	next     alert.Notifier
}

// NewFilterNotifier returns a FilterNotifier that suppresses events matching
// any pattern in the provided list. Comparison is case-insensitive substring
// match against the event's string representation and meta values.
func NewFilterNotifier(patterns []string, next alert.Notifier) *FilterNotifier {
	normalized := make([]string, len(patterns))
	for i, p := range patterns {
		normalized[i] = strings.ToLower(p)
	}
	return &FilterNotifier{patterns: normalized, next: next}
}

// Send checks the event against all patterns. If any pattern matches the
// event string or any meta value the event is silently dropped; otherwise
// it is forwarded to the next notifier.
func (f *FilterNotifier) Send(ev alert.Event) error {
	candidate := strings.ToLower(ev.String())
	for k, v := range ev.Meta {
		candidate += " " + strings.ToLower(k) + "=" + strings.ToLower(v)
	}
	for _, p := range f.patterns {
		if strings.Contains(candidate, p) {
			return nil
		}
	}
	if f.next == nil {
		return nil
	}
	return f.next.Send(ev)
}
