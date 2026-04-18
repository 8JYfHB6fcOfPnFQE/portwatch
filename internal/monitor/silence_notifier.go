package monitor

import (
	"github.com/user/portwatch/internal/alert"
)

// SilenceNotifier wraps another Notifier and drops events for silenced ports.
type SilenceNotifier struct {
	inner  alert.Notifier
	silent *SilenceStore
}

// NewSilenceNotifier returns a SilenceNotifier that filters events via store.
func NewSilenceNotifier(inner alert.Notifier, store *SilenceStore) *SilenceNotifier {
	return &SilenceNotifier{inner: inner, silent: store}
}

// Send forwards the event only when the port is not silenced.
func (n *SilenceNotifier) Send(ev alert.Event) error {
	if n.silent.IsSilenced(ev.Port, ev.Proto) {
		return nil
	}
	return n.inner.Send(ev)
}
