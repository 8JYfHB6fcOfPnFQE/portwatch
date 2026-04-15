package ports

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single open port.
type PortState struct {
	Protocol string
	Port     int
	Address  string
}

// String returns a human-readable representation of the port state.
func (p PortState) String() string {
	return fmt.Sprintf("%s://%s:%d", p.Protocol, p.Address, p.Port)
}

// Scanner scans for open TCP/UDP ports on the local machine.
type Scanner struct {
	Timeout time.Duration
	Ports   []int
}

// NewScanner creates a Scanner with a default timeout and port range.
func NewScanner(timeout time.Duration, ports []int) *Scanner {
	return &Scanner{
		Timeout: timeout,
		Ports:   ports,
	}
}

// Scan checks each port in the configured list and returns those that are open.
func (s *Scanner) Scan() ([]PortState, error) {
	var open []PortState

	for _, port := range s.Ports {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		if err != nil {
			// Port is closed or unreachable — skip
			continue
		}
		conn.Close()

		open = append(open, PortState{
			Protocol: "tcp",
			Port:     port,
			Address:  "127.0.0.1",
		})
	}

	return open, nil
}

// DefaultPortRange returns a commonly monitored range of ports (1–1024).
func DefaultPortRange() []int {
	ports := make([]int, 1024)
	for i := range ports {
		ports[i] = i + 1
	}
	return ports
}
