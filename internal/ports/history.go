package ports

import "sync"

// StateChange represents a transition for a single port between two scans.
type StateChange struct {
	Port  PortState
	OldState string // "open", "closed", or "new"
	NewState string
}

// History tracks the last known port states and computes diffs between scans.
type History struct {
	mu   sync.Mutex
	last map[string]PortState // key: "proto:port"
}

// NewHistory creates an empty History.
func NewHistory() *History {
	return &History{
		last: make(map[string]PortState),
	}
}

// key returns a unique string identifier for a PortState.
func key(p PortState) string {
	return p.Proto + ":" + itoa(p.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 8)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

// Diff compares current scan results against the previous snapshot and returns
// any ports that were opened or closed since the last call.
func (h *History) Diff(current []PortState) []StateChange {
	h.mu.Lock()
	defer h.mu.Unlock()

	currentMap := make(map[string]PortState, len(current))
	for _, p := range current {
		currentMap[key(p)] = p
	}

	var changes []StateChange

	// Detect newly opened ports.
	for k, p := range currentMap {
		if _, existed := h.last[k]; !existed {
			changes = append(changes, StateChange{
				Port:     p,
				OldState: "closed",
				NewState: "open",
			})
		}
	}

	// Detect newly closed ports.
	for k, p := range h.last {
		if _, stillOpen := currentMap[k]; !stillOpen {
			changes = append(changes, StateChange{
				Port:     p,
				OldState: "open",
				NewState: "closed",
			})
		}
	}

	h.last = currentMap
	return changes
}

// Snapshot returns a copy of the current known-open ports.
func (h *History) Snapshot() []PortState {
	h.mu.Lock()
	defer h.mu.Unlock()
	snap := make([]PortState, 0, len(h.last))
	for _, p := range h.last {
		snap = append(snap, p)
	}
	return snap
}
