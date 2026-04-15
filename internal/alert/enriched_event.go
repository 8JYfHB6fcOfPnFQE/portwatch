package alert

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// EnrichedEvent extends Event with process ownership information.
type EnrichedEvent struct {
	Event
	Process ports.ProcessInfo
}

// NewEnrichedEvent constructs an EnrichedEvent from an Event and ProcessInfo.
func NewEnrichedEvent(ev Event, proc ports.ProcessInfo) EnrichedEvent {
	return EnrichedEvent{Event: ev, Process: proc}
}

// String returns a human-readable summary of the enriched event.
func (e EnrichedEvent) String() string {
	return fmt.Sprintf("%s | %s:%d/%s | process=%s",
		e.Timestamp.Format(time.RFC3339),
		e.Kind,
		e.Port.Port,
		e.Port.Proto,
		e.Process.String(),
	)
}

// EnrichingNotifier wraps a Notifier and enriches events with process info
// before forwarding them.
type EnrichingNotifier struct {
	inner   *Notifier
	lookup  func(inode uint64) (ports.ProcessInfo, error)
}

// NewEnrichingNotifier creates an EnrichingNotifier.
func NewEnrichingNotifier(n *Notifier, lookup func(uint64) (ports.ProcessInfo, error)) *EnrichingNotifier {
	return &EnrichingNotifier{inner: n, lookup: lookup}
}

// Send enriches the event with process information and forwards it.
func (en *EnrichingNotifier) Send(ev Event) error {
	var proc ports.ProcessInfo
	if ev.Port.Inode > 0 && en.lookup != nil {
		if info, err := en.lookup(ev.Port.Inode); err == nil {
			proc = info
		}
	}
	enriched := NewEnrichedEvent(ev, proc)
	_, err := fmt.Fprintln(en.inner.writers[0], enriched.String())
	return err
}
