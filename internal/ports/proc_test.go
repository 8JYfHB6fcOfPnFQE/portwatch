package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func writeProcFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "tcp")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write proc file: %v", err)
	}
	return p
}

const procHeader = "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"

func TestParseProcLine_ListenState(t *testing.T) {
	// Port 0x1F90 = 8080
	line := "0: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0"
	ps, ok := parseProcLine(line)
	if !ok {
		t.Fatal("expected ok=true for LISTEN state")
	}
	if ps.Port != 8080 {
		t.Errorf("port: got %d, want 8080", ps.Port)
	}
	if ps.Proto != "tcp" {
		t.Errorf("proto: got %s, want tcp", ps.Proto)
	}
}

func TestParseProcLine_NonListenIgnored(t *testing.T) {
	// State 01 = ESTABLISHED
	line := "1: 00000000:1F90 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0"
	_, ok := parseProcLine(line)
	if ok {
		t.Fatal("expected ok=false for non-LISTEN state")
	}
}

func TestParseProcLine_MalformedLine(t *testing.T) {
	_, ok := parseProcLine("not enough fields")
	if ok {
		t.Fatal("expected ok=false for malformed line")
	}
}

func TestReadProcFile_ParsesListeners(t *testing.T) {
	content := procHeader +
		fmt.Sprintf("  0: 00000000:%04X 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 1 1 0000000000000000 100 0 0 10 0\n", 9090)
	path := writeProcFile(t, content)

	ports, err := readProcFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("got %d ports, want 1", len(ports))
	}
	if ports[0].Port != 9090 {
		t.Errorf("port: got %d, want 9090", ports[0].Port)
	}
}

func TestReadProcFile_MissingFile(t *testing.T) {
	_, err := readProcFile("/nonexistent/path/tcp")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestProcReader_Read_SkipsMissingPaths(t *testing.T) {
	r := &ProcReader{paths: []string{"/nonexistent/tcp", "/nonexistent/tcp6"}}
	ports, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 ports, got %d", len(ports))
	}
}
