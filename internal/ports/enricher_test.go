package ports

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// buildFakeProc creates a minimal /proc-like tree for a single process.
func buildFakeProc(t *testing.T, pid int, comm string, inode uint64) string {
	t.Helper()
	root := t.TempDir()
	pidDir := filepath.Join(root, strconv.Itoa(pid))

	if err := os.MkdirAll(filepath.Join(pidDir, "fd"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte(comm+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Create a symlink that mimics a socket fd.
	socketTarget := "socket:[" + strconv.FormatUint(inode, 10) + "]"
	if err := os.Symlink(socketTarget, filepath.Join(pidDir, "fd", "3")); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestEnricher_Lookup_FindsProcess(t *testing.T) {
	const pid = 1234
	const inode uint64 = 99887766
	root := buildFakeProc(t, pid, "myapp", inode)

	e := NewEnricher(root)
	info, err := e.Lookup(inode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.PID != pid {
		t.Errorf("expected PID %d, got %d", pid, info.PID)
	}
	if info.Name != "myapp" {
		t.Errorf("expected name %q, got %q", "myapp", info.Name)
	}
}

func TestEnricher_Lookup_UnknownInode(t *testing.T) {
	const pid = 42
	root := buildFakeProc(t, pid, "other", 111)

	e := NewEnricher(root)
	info, err := e.Lookup(999) // inode not owned by any process
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.PID != 0 {
		t.Errorf("expected zero PID for unknown inode, got %d", info.PID)
	}
}

func TestProcessInfo_String_WithPID(t *testing.T) {
	p := ProcessInfo{PID: 7, Name: "nginx"}
	got := p.String()
	want := "nginx(pid=7)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestProcessInfo_String_Unknown(t *testing.T) {
	p := ProcessInfo{}
	if got := p.String(); got != "unknown" {
		t.Errorf("String() = %q, want %q", got, "unknown")
	}
}

func TestEnricher_Lookup_BadProcRoot(t *testing.T) {
	e := NewEnricher("/nonexistent/proc")
	_, err := e.Lookup(1)
	if err == nil {
		t.Error("expected error for bad proc root, got nil")
	}
}
