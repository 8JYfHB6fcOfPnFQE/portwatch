package ports

import (
	"fmt"
	"net"
	"time"
)

// PortState represents a single port observed during a scan.
type PortState struct {
	Port  int
	Proto string
	State string
}

// String returns a human-readable representation of the port state.
func (p PortState) String() string {
	return fmt.Sprintf("%s/%d (%s)", p.Proto, p.Port, p.State)
}

// Scanner performs active TCP port scans over a configurable port range.
type Scanner struct {
	From    int
	To      int
	Timeout time.Duration
}

// DefaultPortRange is the range scanned when no explicit range is configured.
const DefaultPortRange = "1-1024"

// NewScanner returns a Scanner with the given port range and a 500 ms dial
// timeout per port.
func NewScanner(from, to int) *Scanner {
	return &Scanner{
		From:    from,
		To:      to,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan iterates over the configured port range and returns every port that
// accepts a TCP connection.
func (s *Scanner) Scan() ([]PortState, error) {
	var open []PortState
	for port := s.From; port <= s.To; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err != nil {
			continue
		}
		conn.Close()
		open = append(open, PortState{
			Port:  port,
			Proto: "tcp",
			State: "LISTEN",
		})
	}
	return open, nil
}
