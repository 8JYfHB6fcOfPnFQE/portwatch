package monitor

import (
	"sync"
	"time"
)

// HealthStatus represents the current health of the monitor daemon.
type HealthStatus struct {
	Healthy     bool      `json:"healthy"`
	LastScan    time.Time `json:"last_scan"`
	LastError   string    `json:"last_error,omitempty"`
	ConsecFails int       `json:"consec_failures"`
	Uptime      string    `json:"uptime"`
	StartedAt   time.Time `json:"started_at"`
}

// HealthTracker tracks the liveness of the monitor loop.
type HealthTracker struct {
	mu          sync.RWMutex
	startedAt   time.Time
	lastScan    time.Time
	lastError   string
	consecFails int
}

// NewHealthTracker creates a new HealthTracker.
func NewHealthTracker() *HealthTracker {
	return &HealthTracker{startedAt: time.Now()}
}

// RecordSuccess marks a successful scan.
func (h *HealthTracker) RecordSuccess() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastScan = time.Now()
	h.lastError = ""
	h.consecFails = 0
}

// RecordFailure marks a failed scan with an error message.
func (h *HealthTracker) RecordFailure(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastScan = time.Now()
	h.lastError = err.Error()
	h.consecFails++
}

// Status returns a snapshot of the current health.
func (h *HealthTracker) Status() HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return HealthStatus{
		Healthy:     h.consecFails == 0,
		LastScan:    h.lastScan,
		LastError:   h.lastError,
		ConsecFails: h.consecFails,
		Uptime:      time.Since(h.startedAt).Round(time.Second).String(),
		StartedAt:   h.startedAt,
	}
}
