package monitor

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// MetricsReporter periodically prints a summary of scan metrics.
type MetricsReporter struct {
	collector *ports.MetricsCollector
	out       io.Writer
	stop      chan struct{}
}

// NewMetricsReporter creates a MetricsReporter writing to w (defaults to os.Stdout).
func NewMetricsReporter(c *ports.MetricsCollector, w io.Writer) *MetricsReporter {
	if w == nil {
		w = os.Stdout
	}
	return &MetricsReporter{collector: c, out: w, stop: make(chan struct{})}
}

// Start begins periodic reporting at the given interval.
func (r *MetricsReporter) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.PrintSummary()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts periodic reporting.
func (r *MetricsReporter) Stop() {
	close(r.stop)
}

// PrintSummary writes a one-line metrics summary to the writer.
func (r *MetricsReporter) PrintSummary() {
	s := r.collector.Summary()
	if s.Total == 0 {
		fmt.Fprintln(r.out, "[metrics] no data collected yet")
		return
	}
	fmt.Fprintf(r.out,
		"[metrics] scans=%d errors=%d avg=%s ports=%d added=%d removed=%d\n",
		s.Total,
		s.Errors,
		roundDuration(s.AvgDuration),
		s.AvgPorts,
		s.TotalAdded,
		s.TotalRemoved,
	)
}

// Reset clears all collected metrics and prints a confirmation message.
func (r *MetricsReporter) Reset() {
	r.collector.Reset()
	fmt.Fprintln(r.out, "[metrics] counters reset")
}

func roundDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}
