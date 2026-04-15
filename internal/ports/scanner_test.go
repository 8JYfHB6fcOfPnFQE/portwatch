package ports_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/ports"
)

// startTestListener opens a TCP listener on a random port and returns the port number and a cleanup func.
func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, cleanup := startTestListener(t)
	defer cleanup()

	scanner := ports.NewScanner(500*time.Millisecond, []int{port})
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(result))
	}
	if result[0].Port != port {
		t.Errorf("expected port %d, got %d", port, result[0].Port)
	}
	if result[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", result[0].Protocol)
	}
}

func TestScan_IgnoresClosedPort(t *testing.T) {
	// Use a port that is very unlikely to be open
	scanner := ports.NewScanner(200*time.Millisecond, []int{19999})
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(result))
	}
}

func TestPortState_String(t *testing.T) {
	ps := ports.PortState{Protocol: "tcp", Address: "127.0.0.1", Port: 8080}
	expected := fmt.Sprintf("tcp://127.0.0.1:8080")
	if ps.String() != expected {
		t.Errorf("expected %q, got %q", expected, ps.String())
	}
}

func TestDefaultPortRange(t *testing.T) {
	r := ports.DefaultPortRange()
	if len(r) != 1024 {
		t.Errorf("expected 1024 ports, got %d", len(r))
	}
	if r[0] != 1 || r[1023] != 1024 {
		t.Errorf("unexpected port range: first=%d last=%d", r[0], r[1023])
	}
}
