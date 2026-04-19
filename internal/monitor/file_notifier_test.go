package monitor_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func makeFileEvent() alert.Event {
	return alert.NewEvent("opened", "tcp", 8080, "127.0.0.1")
}

func TestFileNotifier_InvalidFormat(t *testing.T) {
	_, err := monitor.NewFileNotifier("/tmp/pw_test_invalid.log", "xml", nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestFileNotifier_WritesTextLine(t *testing.T) {
	f, _ := os.CreateTemp("", "pw_text_*.log")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	n, err := monitor.NewFileNotifier(f.Name(), "text", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Send(makeFileEvent()); err != nil {
		t.Fatalf("Send error: %v", err)
	}

	data, _ := os.ReadFile(f.Name())
	line := strings.TrimSpace(string(data))
	if !strings.Contains(line, "opened") || !strings.Contains(line, "tcp") {
		t.Errorf("unexpected text line: %q", line)
	}
}

func TestFileNotifier_WritesJSONLine(t *testing.T) {
	f, _ := os.CreateTemp("", "pw_json_*.log")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	n, err := monitor.NewFileNotifier(f.Name(), "json", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Send(makeFileEvent()); err != nil {
		t.Fatalf("Send error: %v", err)
	}

	data, _ := os.ReadFile(f.Name())
	var rec struct {
		Kind  string `json:"kind"`
		Proto string `json:"proto"`
		Port  int    `json:"port"`
		Time  string `json:"time"`
	}{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec.Kind != "opened" || rec.Proto != "tcp" || rec.Port != 8080 {
		t.Errorf("unexpected record: %+v", rec)
	}
	if _, err := time.Parse(time.RFC3339, rec.Time); err != nil {
		t.Errorf("bad time format: %v", err)
	}
}

func TestFileNotifier_ForwardsToNext(t *testing.T) {
	f, _ := os.CreateTemp("", "pw_fwd_*.log")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error { got = e; return nil })

	n, _ := monitor.NewFileNotifier(f.Name(), "text", next)
	e := makeFileEvent()
	n.Send(e)

	if got.Port != e.Port {
		t.Errorf("next not called with correct event")
	}
}
