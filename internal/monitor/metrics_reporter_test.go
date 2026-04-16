package monitor_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func makeCollectorWithData(t *testing.T) *ports.MetricsCollector {
	t.Helper()
	c := ports.NewMetricsCollector(10)
	c.Record(ports.ScanMetrics{
		Timestamp:  time.Now(),
		Duration:   120 * time.Millisecond,
		PortsFound: 5,
		Added:      1,
		Removed:    0,
		Error:      false,
	})
	c.Record(ports.ScanMetrics{
		Timestamp:  time.Now(),
		Duration:   95 * time.Millisecond,
		PortsFound: 6,
		Added:      0,
		Removed:    1,
		Error:      false,
	})
	return c
}

func TestMetricsReporter_PrintSummary(t *testing.T) {
	c := makeCollectorWithData(t)
	var buf bytes.Buffer
	r := monitor.NewMetricsReporter(c, &buf)
	r.PrintSummary()
	out := buf.String()
	if !strings.Contains(out, "scans") {
		t.Errorf("expected 'scans' in output, got: %s", out)
	}
	if !strings.Contains(out, "avg") {
		t.Errorf("expected 'avg' in output, got: %s", out)
	}
}

func TestMetricsReporter_PrintSummary_Empty(t *testing.T) {
	c := ports.NewMetricsCollector(10)
	var buf bytes.Buffer
	r := monitor.NewMetricsReporter(c, &buf)
	r.PrintSummary()
	out := buf.String()
	if !strings.Contains(out, "no data") {
		t.Errorf("expected 'no data' in output, got: %s", out)
	}
}

func TestMetricsReporter_StartStop(t *testing.T) {
	c := makeCollectorWithData(t)
	var buf bytes.Buffer
	r := monitor.NewMetricsReporter(c, &buf)
	r.Start(50 * time.Millisecond)
	time.Sleep(130 * time.Millisecond)
	r.Stop()
	out := buf.String()
	if !strings.Contains(out, "scans") {
		t.Errorf("expected periodic output, got: %s", out)
	}
}
