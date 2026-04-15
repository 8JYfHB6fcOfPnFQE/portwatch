// Package ports provides port scanning, filtering, and history tracking.
package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProcReader reads open ports directly from /proc/net/tcp and /proc/net/tcp6.
type ProcReader struct {
	paths []string
}

// NewProcReader returns a ProcReader targeting the default proc paths.
func NewProcReader() *ProcReader {
	return &ProcReader{
		paths: []string{"/proc/net/tcp", "/proc/net/tcp6"},
	}
}

// Read returns all LISTEN-state ports found in /proc/net/tcp[6].
func (r *ProcReader) Read() ([]PortState, error) {
	var result []PortState
	for _, path := range r.paths {
		ports, err := readProcFile(path)
		if err != nil {
			// Non-fatal: file may not exist on non-Linux systems.
			continue
		}
		result = append(result, ports...)
	}
	return result, nil
}

// readProcFile parses a single /proc/net/tcp or /proc/net/tcp6 file.
func readProcFile(path string) ([]PortState, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var ports []PortState
	scanner := bufio.NewScanner(f)
	// Skip header line.
	scanner.Scan()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ps, ok := parseProcLine(line)
		if ok {
			ports = append(ports, ps)
		}
	}
	return ports, scanner.Err()
}

// parseProcLine extracts a PortState from a /proc/net/tcp line.
// Returns false if the entry is not in LISTEN state (0A).
func parseProcLine(line string) (PortState, bool) {
	fields := strings.Fields(line)
	// Need at least: sl local_addr rem_addr st ...
	if len(fields) < 4 {
		return PortState{}, false
	}
	// State field: 0A = TCP_LISTEN
	if strings.ToUpper(fields[3]) != "0A" {
		return PortState{}, false
	}
	// local_addr is hex "IPADDR:PORT"
	parts := strings.SplitN(fields[1], ":", 2)
	if len(parts) != 2 {
		return PortState{}, false
	}
	portHex := parts[1]
	portNum, err := strconv.ParseUint(portHex, 16, 16)
	if err != nil {
		return PortState{}, false
	}
	return PortState{Port: int(portNum), Proto: "tcp", State: "LISTEN"}, true
}
