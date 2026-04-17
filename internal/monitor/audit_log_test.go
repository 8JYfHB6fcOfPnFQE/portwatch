package monitor

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
)

func makeAuditEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Kind:      kind,
		State:     ports.PortState{Proto: proto, Port: port},
	}
}

func TestAuditLog_RecordWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	al := NewAuditLogWriter(&buf)

	ev := makeAuditEvent("opened", "tcp", 8080)
	if err := al.Record(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var entry AuditEntry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Port != 8080 || entry.Proto != "tcp" || entry.Kind != "opened" {
		t.Errorf("unexpected entry: %+v", entry)
	}
	if entry.Action != "alert" {
		t.Errorf("expected action=alert, got %s", entry.Action)
	}
}

func TestAuditLog_MultipleRecords(t *testing.T) {
	var buf bytes.Buffer
	al := NewAuditLogWriter(&buf)

	for _, port := range []uint16{22, 80, 443} {
		if err := al.Record(makeAuditEvent("opened", "tcp", port)); err != nil {
			t.Fatalf("record error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestAuditNotifier_ForwardsEvent(t *testing.T) {
	var buf bytes.Buffer
	var received []alert.Event

	stub := alert.NotifierFunc(func(ev alert.Event) error {
		received = append(received, ev)
		return nil
	})

	an := NewAuditNotifier(stub, NewAuditLogWriter(&buf))
	ev := makeAuditEvent("closed", "udp", 53)
	if err := an.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 forwarded event, got %d", len(received))
	}
	if buf.Len() == 0 {
		t.Error("expected audit log to be written")
	}
}
