package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// AuditEntry represents a single persisted alert event.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      string    `json:"kind"`
	Proto     string    `json:"proto"`
	Port      uint16    `json:"port"`
	Action    string    `json:"action"`
}

// AuditLog writes alert events to a file in newline-delimited JSON.
type AuditLog struct {
	w io.Writer
}

// NewAuditLog opens (or creates) the file at path for append-only writing.
func NewAuditLog(path string) (*AuditLog, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit log: open %s: %w", path, err)
	}
	return &AuditLog{w: f}, nil
}

// NewAuditLogWriter creates an AuditLog writing to an arbitrary writer (useful for tests).
func NewAuditLogWriter(w io.Writer) *AuditLog {
	return &AuditLog{w: w}
}

// Record encodes the event as a JSON line.
func (a *AuditLog) Record(ev alert.Event) error {
	entry := AuditEntry{
		Timestamp: ev.Timestamp,
		Kind:      ev.Kind,
		Proto:     ev.State.Proto,
		Port:      ev.State.Port,
		Action:    "alert",
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit log: marshal: %w", err)
	}
	_, err = fmt.Fprintf(a.w, "%s\n", b)
	return err
}
