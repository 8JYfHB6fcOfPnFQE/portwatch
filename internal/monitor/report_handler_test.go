package monitor

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

type mockReportScanner struct {
	states []ports.PortState
	err    error
}

func (m *mockReportScanner) Scan() ([]ports.PortState, error) { return m.states, m.err }

type mockReportEnricher struct{}

func (m *mockReportEnricher) Lookup(inode uint64) (ports.ProcessInfo, error) {
	if inode == 42 {
		return ports.ProcessInfo{PID: 99, Name: "nginx"}, nil
	}
	return ports.ProcessInfo{}, errors.New("not found")
}

func TestReportHandler_Print_Success(t *testing.T) {
	scanner := &mockReportScanner{
		states: []ports.PortState{
			{Proto: "tcp", Port: 80, State: "LISTEN", Inode: 42},
		},
	}
	var buf strings.Builder
	h := NewReportHandler(scanner, &mockReportEnricher{}, &buf)
	if err := h.Print(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Errorf("expected port 80 in output, got: %s", out)
	}
	if !strings.Contains(out, "nginx") {
		t.Errorf("expected process name in output, got: %s", out)
	}
}

func TestReportHandler_Print_ScanError(t *testing.T) {
	scanner := &mockReportScanner{err: errors.New("scan failed")}
	var buf strings.Builder
	h := NewReportHandler(scanner, nil, &buf)
	if err := h.Print(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReportHandler_Print_NilEnricher(t *testing.T) {
	scanner := &mockReportScanner{
		states: []ports.PortState{
			{Proto: "udp", Port: 53, State: "LISTEN", Inode: 7},
		},
	}
	var buf strings.Builder
	h := NewReportHandler(scanner, nil, &buf)
	if err := h.Print(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "unknown") {
		t.Errorf("expected 'unknown' process when enricher is nil")
	}
}

func TestReportHandler_DefaultsToStdout(t *testing.T) {
	h := NewReportHandler(&mockReportScanner{}, nil, nil)
	if h.out == nil {
		t.Error("expected non-nil writer")
	}
}
