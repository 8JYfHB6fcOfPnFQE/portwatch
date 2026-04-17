package monitor

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// AuditNotifier wraps an inner Notifier and additionally records every event
// to an AuditLog before forwarding it.
type AuditNotifier struct {
	inner alert.Notifier
	log   *AuditLog
}

// NewAuditNotifier creates an AuditNotifier.
func NewAuditNotifier(inner alert.Notifier, log *AuditLog) *AuditNotifier {
	return &AuditNotifier{inner: inner, log: log}
}

// Send records the event then delegates to the inner Notifier.
func (a *AuditNotifier) Send(ev alert.Event) error {
	if err := a.log.Record(ev); err != nil {
		// Non-fatal: log the error but still forward the event.
		fmt.Printf("portwatch: audit log write error: %v\n", err)
	}
	return a.inner.Send(ev)
}
