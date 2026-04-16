package monitor

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

type fakeScanner struct {
	states []ports.PortState
	err    error
}

func (f *fakeScanner) Scan() ([]ports.PortState, error) {
	return f.states, f.err
}

func TestBaselineHandler_Capture(t *testing.T) {
	var buf bytes.Buffer
	store := ports.NewBaselineStore(filepath.Join(t.TempDir(), "bl.json"))
	scanner := &fakeScanner{states: []ports.PortState{{Port: 22, Proto: "tcp"}}}
	h := NewBaselineHandler(store, scanner, &buf)
	if err := h.Capture(); err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if !strings.Contains(buf.String(), "1 ports saved") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestBaselineHandler_Capture_ScanError(t *testing.T) {
	var buf bytes.Buffer
	store := ports.NewBaselineStore(filepath.Join(t.TempDir(), "bl.json"))
	scanner := &fakeScanner{err: errors.New("scan fail")}
	h := NewBaselineHandler(store, scanner, &buf)
	if err := h.Capture(); err == nil {
		t.Error("expected error")
	}
}

func TestBaselineHandler_Show_NoBaseline(t *testing.T) {
	var buf bytes.Buffer
	store := ports.NewBaselineStore("/nonexistent/bl.json")
	h := NewBaselineHandler(store, &fakeScanner{}, &buf)
	if err := h.Show(); err != nil {
		t.Fatalf("Show: %v", err)
	}
	if !strings.Contains(buf.String(), "no baseline found") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestBaselineHandler_Check_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	store := ports.NewBaselineStore(filepath.Join(t.TempDir(), "bl.json"))
	states := []ports.PortState{{Port: 80, Proto: "tcp"}}
	scanner := &fakeScanner{states: states}
	h := NewBaselineHandler(store, scanner, &buf)
	_ = h.Capture()
	buf.Reset()
	if err := h.Check(); err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestBaselineHandler_Check_DetectsAdded(t *testing.T) {
	var buf bytes.Buffer
	store := ports.NewBaselineStore(filepath.Join(t.TempDir(), "bl.json"))
	original := &fakeScanner{states: []ports.PortState{{Port: 80, Proto: "tcp"}}}
	h := NewBaselineHandler(store, original, &buf)
	_ = h.Capture()
	h.scanner = &fakeScanner{states: []ports.PortState{
		{Port: 80, Proto: "tcp"},
		{Port: 9090, Proto: "tcp"},
	}}
	buf.Reset()
	_ = h.Check()
	if !strings.Contains(buf.String(), "ADDED") {
		t.Errorf("expected ADDED in output: %s", buf.String())
	}
}
