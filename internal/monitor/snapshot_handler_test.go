package monitor

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func tempSnapPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestSnapshotHandler_Capture_And_Show(t *testing.T) {
	ln := startTCPListener(t)
	port := listenerPort(t, ln)

	scanner := ports.NewScanner(port, port)
	store := ports.NewSnapshotStore(tempSnapPath(t))
	h := NewSnapshotHandler(store, scanner)
	h.out = &bytes.Buffer{}

	if err := h.Capture(context.Background()); err != nil {
		t.Fatalf("Capture: %v", err)
	}

	var buf bytes.Buffer
	h.out = &buf
	if err := h.Show(); err != nil {
		t.Fatalf("Show: %v", err)
	}

	if !strings.Contains(buf.String(), "PROTO") {
		t.Errorf("expected header in output, got: %s", buf.String())
	}
}

func TestSnapshotHandler_Show_NoSnapshot(t *testing.T) {
	store := ports.NewSnapshotStore(filepath.Join(t.TempDir(), "missing.json"))
	scanner := ports.NewScanner(9900, 9900)
	h := NewSnapshotHandler(store, scanner)

	var buf bytes.Buffer
	h.out = &buf

	if err := h.Show(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no snapshot") {
		t.Errorf("expected 'no snapshot' message, got: %s", buf.String())
	}
}

func TestSnapshotHandler_Capture_ScanError(t *testing.T) {
	store := ports.NewSnapshotStore(filepath.Join(t.TempDir(), "snap.json"))
	// scanner with invalid range to force error
	scanner := ports.NewScanner(0, 0)
	h := NewSnapshotHandler(store, scanner)

	// should not panic; error is wrapped
	_ = h.Capture(context.Background())
}

func listenerPort(t *testing.T, ln interface{ Addr() interface{ String() string } }) int {
	t.Helper()
	addr := ln.Addr().String()
	var port int
	_, _ = fmt.Sscanf(addr[strings.LastIndex(addr, ":")+1:], "%d", &port)
	return port
}

func init() {
	_ = os.Getenv // suppress unused import
}
