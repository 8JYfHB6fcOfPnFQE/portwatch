package monitor

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// MetricsHandler wraps a scan function and records ScanMetrics for each cycle.
type MetricsHandler struct {
	collector *ports.MetricsCollector
	scan      func(ctx context.Context) ([]ports.PortState, error)
	prev      []ports.PortState
}

// NewMetricsHandler creates a MetricsHandler using the provided collector and scan func.
func NewMetricsHandler(collector *ports.MetricsCollector, scan func(ctx context.Context) ([]ports.PortState, error)) *MetricsHandler {
	return &MetricsHandler{collector: collector, scan: scan}
}

// Run executes one scan cycle, records metrics, and returns current port states.
func (h *MetricsHandler) Run(ctx context.Context) ([]ports.PortState, error) {
	start := time.Now()
	states, err := h.scan(ctx)
	duration := time.Since(start)

	sm := ports.ScanMetrics{
		ScanDuration: duration,
		ScanError:    err != nil,
		Timestamp:    time.Now(),
	}

	if err == nil {
		sm.PortsFound = len(states)
		added, removed := diffCounts(h.prev, states)
		sm.PortsAdded = added
		sm.PortsRemoved = removed
		h.prev = states
	}

	h.collector.Record(sm)
	return states, err
}

// diffCounts returns the number of added and removed ports between two snapshots.
func diffCounts(prev, curr []ports.PortState) (added, removed int) {
	prevSet := make(map[string]struct{}, len(prev))
	for _, p := range prev {
		prevSet[p.String()] = struct{}{}
	}
	currSet := make(map[string]struct{}, len(curr))
	for _, c := range curr {
		currSet[c.String()] = struct{}{}
	}
	for k := range currSet {
		if _, ok := prevSet[k]; !ok {
			added++
		}
	}
	for k := range prevSet {
		if _, ok := currSet[k]; !ok {
			removed++
		}
	}
	return
}
