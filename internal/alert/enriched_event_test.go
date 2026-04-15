package alert

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makeTestEvent() Event {
	return NewEvent("opened", ports.PortState{
		Port:  8080,
		Proto: "tcp",
		Inode: 42,
	})
}

func TestEnrichedEvent_String_ContainsProcess(t *testing.T) {
	ev := makeTestEvent()
	proc := ports.ProcessInfo{PID: 99, Name: "server"}
	ee := NewEnrichedEvent(ev, proc)

	got := ee.String()
	if !strings.Contains(got, "server(pid=99)") {
		t.Errorf("String() missing process info, got: %s", got)
	}
	if !strings.Contains(got, "8080") {
		t.Errorf("String() missing port, got: %s", got)
	}
	if !strings.Contains(got, "tcp") {
		t.Errorf("String() missing proto, got: %s", got)
	}
}

func TestEnrichedEvent_String_UnknownProcess(t *testing.T) {
	ev := makeTestEvent()
	ee := NewEnrichedEvent(ev, ports.ProcessInfo{})
	if !strings.Contains(ee.String(), "unknown") {
		t.Errorf("expected 'unknown' for zero ProcessInfo, got: %s", ee.String())
	}
}

func TestEnrichingNotifier_Send_WithLookup(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	lookup := func(inode uint64) (ports.ProcessInfo, error) {
		return ports.ProcessInfo{PID: 5, Name: "daemon"}, nil
	}
	en := NewEnrichingNotifier(n, lookup)

	ev := makeTestEvent()
	ev.Timestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := en.Send(ev); err != nil {
		t.Fatalf("Send() error: %v", err)
	}
	line := buf.String()
	if !strings.Contains(line, "daemon(pid=5)") {
		t.Errorf("expected enriched process info in output, got: %s", line)
	}
}

func TestEnrichingNotifier_Send_LookupError(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	lookup := func(inode uint64) (ports.ProcessInfo, error) {
		return ports.ProcessInfo{}, errors.New("lookup failed")
	}
	en := NewEnrichingNotifier(n, lookup)

	ev := makeTestEvent()
	if err := en.Send(ev); err != nil {
		t.Fatalf("Send() should not error on lookup failure, got: %v", err)
	}
	if !strings.Contains(buf.String(), "unknown") {
		t.Errorf("expected 'unknown' when lookup fails, got: %s", buf.String())
	}
}
