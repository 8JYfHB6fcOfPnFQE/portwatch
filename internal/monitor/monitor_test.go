package monitor_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func startTCPListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestMonitor_DetectsOpenedPort(t *testing.T) {
	port, closeListener := startTCPListener(t)
	defer closeListener()

	scanner := ports.NewScanner(ports.PortRange{Start: port, End: port})
	m := monitor.New(scanner, nil, 50*time.Millisecond)
	m.Start()
	defer m.Stop()

	select {
	case change := <-m.Changes:
		if change.Change != monitor.ChangeOpened {
			t.Errorf("expected ChangeOpened, got %s", change.Change)
		}
		if change.Port.Port != port {
			t.Errorf("expected port %d, got %d", port, change.Port.Port)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for ChangeOpened event")
	}
}

func TestMonitor_DetectsClosedPort(t *testing.T) {
	port, closeListener := startTCPListener(t)

	scanner := ports.NewScanner(ports.PortRange{Start: port, End: port})
	m := monitor.New(scanner, nil, 50*time.Millisecond)
	m.Start()
	defer m.Stop()

	// Wait for the open event first.
	select {
	case <-m.Changes:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for initial open event")
	}

	closeListener()

	select {
	case change := <-m.Changes:
		if change.Change != monitor.ChangeClosed {
			t.Errorf("expected ChangeClosed, got %s", change.Change)
		}
		if change.Port.Port != port {
			t.Errorf("expected port %d, got %d", port, change.Port.Port)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for ChangeClosed event")
	}
}

func TestMonitor_Stop(t *testing.T) {
	scanner := ports.NewScanner(ports.DefaultPortRange())
	m := monitor.New(scanner, nil, 10*time.Second)
	m.Start()
	done := make(chan struct{})
	go func() {
		m.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Stop() did not return in time")
	}
}
