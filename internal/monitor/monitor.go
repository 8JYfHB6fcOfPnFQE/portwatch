package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

// ChangeType describes what kind of port change occurred.
type ChangeType string

const (
	ChangeOpened ChangeType = "opened"
	ChangeClosed ChangeType = "closed"
)

// PortChange represents a detected change in port state.
type PortChange struct {
	Port   ports.PortState
	Change ChangeType
}

// Monitor watches ports at a given interval and emits changes.
type Monitor struct {
	scanner  *ports.Scanner
	matcher  *rules.Matcher
	interval time.Duration
	previous map[string]ports.PortState
	Changes  chan PortChange
	quit     chan struct{}
}

// New creates a Monitor with the given scanner, matcher, and poll interval.
func New(scanner *ports.Scanner, matcher *rules.Matcher, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  scanner,
		matcher:  matcher,
		interval: interval,
		previous: make(map[string]ports.PortState),
		Changes:  make(chan PortChange, 64),
		quit:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (m *Monitor) Start() {
	go m.loop()
}

// Stop signals the monitor to cease polling.
func (m *Monitor) Stop() {
	close(m.quit)
}

func (m *Monitor) loop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.poll()
		case <-m.quit:
			return
		}
	}
}

func (m *Monitor) poll() {
	current, err := m.scanner.Scan()
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return
	}

	currentMap := make(map[string]ports.PortState, len(current))
	for _, ps := range current {
		key := ps.String()
		currentMap[key] = ps
		if _, existed := m.previous[key]; !existed {
			if m.matcher != nil {
				m.matcher.Match(ps)
			}
			m.Changes <- PortChange{Port: ps, Change: ChangeOpened}
		}
	}

	for key, ps := range m.previous {
		if _, stillOpen := currentMap[key]; !stillOpen {
			m.Changes <- PortChange{Port: ps, Change: ChangeClosed}
		}
	}

	m.previous = currentMap
}
