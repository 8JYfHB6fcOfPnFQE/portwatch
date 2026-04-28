package monitor

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// EnvelopeNotifier wraps each event in a metadata envelope before forwarding.
// It stamps events with a sequence number, a send timestamp, and an optional
// source identifier so downstream consumers can correlate and order alerts.
type EnvelopeNotifier struct {
	source string
	next   alert.Notifier
	seq    uint64
}

// NewEnvelopeNotifier returns an EnvelopeNotifier that tags every event with
// envelope metadata and forwards it to next. source is an arbitrary label
// (e.g. hostname or cluster name) embedded in every envelope.
func NewEnvelopeNotifier(source string, next alert.Notifier) *EnvelopeNotifier {
	return &EnvelopeNotifier{source: source, next: next}
}

// Send stamps the event with envelope metadata and forwards it.
func (e *EnvelopeNotifier) Send(ev alert.Event) error {
	e.seq++

	if ev.Meta == nil {
		ev.Meta = make(map[string]string)
	}
	ev.Meta["envelope.seq"] = fmt.Sprintf("%d", e.seq)
	ev.Meta["envelope.sent_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	if e.source != "" {
		ev.Meta["envelope.source"] = e.source
	}

	if e.next == nil {
		return nil
	}
	return e.next.Send(ev)
}
