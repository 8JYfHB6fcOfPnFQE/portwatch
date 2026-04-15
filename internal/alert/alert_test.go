package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewNotifier_DefaultsToStdout(t *testing.T) {
	n := NewNotifier()
	if len(n.writers) != 1 {
		t.Fatalf("expected 1 writer, got %d", len(n.writers))
	}
}

func TestNotifier_Send_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	e := Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     LevelAlert,
		Port:      8080,
		Proto:     "tcp",
		Message:   "unexpected port opened",
	}
	n.Send(e)

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "unexpected port opened") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestNotifier_Send_MultipleWriters(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	n := NewNotifier(&buf1, &buf2)

	n.Send(NewEvent(LevelWarn, 22, "tcp", "ssh port active"))

	if buf1.Len() == 0 {
		t.Error("expected buf1 to have content")
	}
	if buf2.Len() == 0 {
		t.Error("expected buf2 to have content")
	}
	if buf1.String() != buf2.String() {
		t.Error("expected both writers to receive identical output")
	}
}

func TestNewEvent_SetsFields(t *testing.T) {
	before := time.Now()
	e := NewEvent(LevelInfo, 443, "tcp", "port allowed")
	after := time.Now()

	if e.Level != LevelInfo {
		t.Errorf("expected LevelInfo, got %s", e.Level)
	}
	if e.Port != 443 {
		t.Errorf("expected port 443, got %d", e.Port)
	}
	if e.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", e.Proto)
	}
	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Error("timestamp out of expected range")
	}
}
