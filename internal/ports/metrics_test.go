package ports

import (
	"testing"
	"time"
)

func TestMetricsCollector_RecordAndLatest(t *testing.T) {
	mc := NewMetricsCollector(10)
	_, ok := mc.Latest()
	if ok {
		t.Fatal("expected no latest on empty collector")
	}

	sm := ScanMetrics{PortsFound: 5, PortsAdded: 2, ScanDuration: 10 * time.Millisecond}
	mc.Record(sm)

	got, ok := mc.Latest()
	if !ok {
		t.Fatal("expected latest after record")
	}
	if got.PortsFound != 5 || got.PortsAdded != 2 {
		t.Errorf("unexpected latest: %+v", got)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestMetricsCollector_CapHistory(t *testing.T) {
	mc := NewMetricsCollector(3)
	for i := 0; i < 5; i++ {
		mc.Record(ScanMetrics{PortsAdded: i})
	}
	mc.mu.Lock()
	l := len(mc.history)
	mc.mu.Unlock()
	if l != 3 {
		t.Errorf("expected 3 entries, got %d", l)
	}
}

func TestMetricsCollector_Summary(t *testing.T) {
	mc := NewMetricsCollector(10)
	mc.Record(ScanMetrics{PortsAdded: 2, PortsRemoved: 1, ScanDuration: 20 * time.Millisecond})
	mc.Record(ScanMetrics{PortsAdded: 1, ScanError: true, ScanDuration: 40 * time.Millisecond})

	s := mc.Summary()
	if s.Count != 2 {
		t.Errorf("expected count 2, got %d", s.Count)
	}
	if s.TotalAdded != 3 {
		t.Errorf("expected TotalAdded 3, got %d", s.TotalAdded)
	}
	if s.TotalRemoved != 1 {
		t.Errorf("expected TotalRemoved 1, got %d", s.TotalRemoved)
	}
	if s.TotalErrors != 1 {
		t.Errorf("expected TotalErrors 1, got %d", s.TotalErrors)
	}
	if s.AvgDuration != 30*time.Millisecond {
		t.Errorf("expected AvgDuration 30ms, got %v", s.AvgDuration)
	}
}

func TestMetricsCollector_DefaultMaxLen(t *testing.T) {
	mc := NewMetricsCollector(0)
	if mc.maxLen != 100 {
		t.Errorf("expected default maxLen 100, got %d", mc.maxLen)
	}
}
