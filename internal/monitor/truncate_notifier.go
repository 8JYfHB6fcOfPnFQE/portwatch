package monitor

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// TruncateNotifier trims string fields in an event that exceed a configured
// maximum length. This is useful when forwarding to backends that enforce
// payload size limits (e.g. SMS gateways, certain webhook endpoints).
//
// Fields affected: Message, Addr, and any meta values.
type TruncateNotifier struct {
	next      alert.Notifier
	maxLen    int
	ellipsis  string
}

// NewTruncateNotifier returns a TruncateNotifier that truncates string fields
// to maxLen runes. If maxLen <= 0 the notifier forwards events unchanged.
func NewTruncateNotifier(maxLen int, next alert.Notifier) *TruncateNotifier {
	return &TruncateNotifier{
		next:     next,
		maxLen:   maxLen,
		ellipsis: "...",
	}
}

func (t *TruncateNotifier) Send(e alert.Event) error {
	if t.maxLen <= 0 {
		if t.next != nil {
			return t.next.Send(e)
		}
		return nil
	}

	e.Message = t.truncate(e.Message)
	e.Addr = t.truncate(e.Addr)

	if len(e.Meta) > 0 {
		copied := make(map[string]string, len(e.Meta))
		for k, v := range e.Meta {
			copied[k] = t.truncate(v)
		}
		e.Meta = copied
	}

	if t.next != nil {
		return t.next.Send(e)
	}
	return nil
}

func (t *TruncateNotifier) truncate(s string) string {
	runes := []rune(s)
	if len(runes) <= t.maxLen {
		return s
	}
	cutAt := t.maxLen - len([]rune(t.ellipsis))
	if cutAt < 0 {
		cutAt = 0
	}
	return strings.TrimSpace(string(runes[:cutAt])) + t.ellipsis
}
