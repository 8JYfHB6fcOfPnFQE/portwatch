package monitor

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SuppressHandler filters out port change events that occur within a quiet
// window after startup, avoiding alert storms on daemon restart.
type SuppressHandler struct {
	mu        sync.Mutex
	window    time.Duration
	startedAt time.Time
	suppressed int
}

// NewSuppressHandler creates a SuppressHandler that silences alerts for the
// given duration after creation.
func NewSuppressHandler(window time.Duration) *SuppressHandler {
	return &SuppressHandler{
		window:    window,
		startedAt: time.Now(),
	}
}

// IsSuppressed returns true if the handler is still within the quiet window.
func (s *SuppressHandler) IsSuppressed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return time.Since(s.startedAt) < s.window
}

// FilterDiff removes diff entries when inside the suppression window and
// records how many events were dropped. Returns the (possibly empty) diff.
func (s *SuppressHandler) FilterDiff(diff []ports.PortState) []ports.PortState {
	if !s.IsSuppressed() {
		return diff
	}
	s.mu.Lock()
	s.suppressed += len(diff)
	s.mu.Unlock()
	return nil
}

// Suppressed returns the total number of suppressed events so far.
func (s *SuppressHandler) Suppressed() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.suppressed
}

// Reset restarts the suppression window from now.
func (s *SuppressHandler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.startedAt = time.Now()
	s.suppressed = 0
}
