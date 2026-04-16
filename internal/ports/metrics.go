package ports

import (
	"sync"
	"time"
)

// ScanMetrics holds counters and timing for a single scan cycle.
type ScanMetrics struct {
	ScanDuration time.Duration
	PortsFound   int
	PortsAdded   int
	PortsRemoved int
	ScanError    bool
	Timestamp    time.Time
}

// MetricsCollector accumulates scan metrics over time.
type MetricsCollector struct {
	mu      sync.Mutex
	history []ScanMetrics
	maxLen  int
}

// NewMetricsCollector creates a MetricsCollector that retains up to maxLen entries.
func NewMetricsCollector(maxLen int) *MetricsCollector {
	if maxLen <= 0 {
		maxLen = 100
	}
	return &MetricsCollector{maxLen: maxLen}
}

// Record appends a ScanMetrics snapshot.
func (m *MetricsCollector) Record(sm ScanMetrics) {
	if sm.Timestamp.IsZero() {
		sm.Timestamp = time.Now()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.history) >= m.maxLen {
		m.history = m.history[1:]
	}
	m.history = append(m.history, sm)
}

// Latest returns the most recent ScanMetrics, or false if none recorded.
func (m *MetricsCollector) Latest() (ScanMetrics, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.history) == 0 {
		return ScanMetrics{}, false
	}
	return m.history[len(m.history)-1], true
}

// Summary returns aggregate stats across all retained metrics.
func (m *MetricsCollector) Summary() MetricsSummary {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := MetricsSummary{Count: len(m.history)}
	for _, sm := range m.history {
		s.TotalAdded += sm.PortsAdded
		s.TotalRemoved += sm.PortsRemoved
		s.TotalErrors += boolInt(sm.ScanError)
		s.TotalDuration += sm.ScanDuration
	}
	if s.Count > 0 {
		s.AvgDuration = s.TotalDuration / time.Duration(s.Count)
	}
	return s
}

// MetricsSummary aggregates stats across multiple scans.
type MetricsSummary struct {
	Count         int
	TotalAdded    int
	TotalRemoved  int
	TotalErrors   int
	TotalDuration time.Duration
	AvgDuration   time.Duration
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
