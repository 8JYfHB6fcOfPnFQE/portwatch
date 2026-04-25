package monitor

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// TransformFunc mutates an event before forwarding it.
type TransformFunc func(e alert.Event) alert.Event

// TransformNotifier applies a sequence of TransformFuncs to each event
// before passing it downstream. This allows lightweight, composable
// event mutation without requiring a full notifier wrapper per change.
type TransformNotifier struct {
	next       alert.Notifier
	transforms []TransformFunc
}

// NewTransformNotifier returns a TransformNotifier that applies fns in
// order and forwards the result to next. If next is nil the event is
// silently dropped after transformation.
func NewTransformNotifier(next alert.Notifier, fns ...TransformFunc) *TransformNotifier {
	return &TransformNotifier{next: next, transforms: fns}
}

// Send applies all registered transforms then forwards to next.
func (t *TransformNotifier) Send(e alert.Event) error {
	for _, fn := range t.transforms {
		e = fn(e)
	}
	if t.next == nil {
		return nil
	}
	return t.next.Send(e)
}

// --- built-in transform helpers ---

// UpperCaseAction returns a TransformFunc that upper-cases Event.Action.
func UpperCaseAction() TransformFunc {
	return func(e alert.Event) alert.Event {
		e.Action = strings.ToUpper(e.Action)
		return e
	}
}

// SetMeta returns a TransformFunc that sets a single metadata key.
func SetMeta(key, value string) TransformFunc {
	return func(e alert.Event) alert.Event {
		if e.Meta == nil {
			e.Meta = make(map[string]string)
		}
		e.Meta[key] = value
		return e
	}
}

// RedactAddr returns a TransformFunc that blanks the remote address field
// so it is not forwarded to external notifiers.
func RedactAddr() TransformFunc {
	return func(e alert.Event) alert.Event {
		e.Addr = ""
		return e
	}
}
