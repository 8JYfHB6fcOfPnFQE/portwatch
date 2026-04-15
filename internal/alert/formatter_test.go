package alert

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestFormattedNotifier_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	fn := NewFormattedNotifier(FormatText, &buf)

	e := Event{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		Port:      3306,
		Proto:     "tcp",
		Message:   "database port exposed",
	}
	fn.Send(e)

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("text output missing level, got: %s", out)
	}
	if !strings.Contains(out, "3306") {
		t.Errorf("text output missing port, got: %s", out)
	}
}

func TestFormattedNotifier_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	fn := NewFormattedNotifier(FormatJSON, &buf)

	e := Event{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Level:     LevelAlert,
		Port:      9200,
		Proto:     "tcp",
		Message:   "elasticsearch exposed",
	}
	fn.Send(e)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if result["level"] != string(LevelAlert) {
		t.Errorf("expected level ALERT, got %v", result["level"])
	}
	if result["port"] != float64(9200) {
		t.Errorf("expected port 9200, got %v", result["port"])
	}
	if result["message"] != "elasticsearch exposed" {
		t.Errorf("unexpected message: %v", result["message"])
	}
	if result["timestamp"] != "2024-06-01T12:00:00Z" {
		t.Errorf("unexpected timestamp: %v", result["timestamp"])
	}
}
