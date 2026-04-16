package monitor

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// MetricsReporter prints scan metrics summaries to a writer.
type MetricsReporter struct {
	collector *ports.MetricsCollector
	out       io.Writer
}

// NewMetricsReporter creates a reporter writing to out.
func NewMetricsReporter(collector *ports.MetricsCollector, out io.Writer) *MetricsReporter {
	return &MetricsReporter{collector: collector, out: out}
}

// PrintSummary writes aggregate stats to the writer.
func (r *MetricsReporter) PrintSummary() {
	s := r.collector.Summary()
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "=== portwatch scan metrics ===")
	fmt.Fprintf(w, "Scans completed:\t%d\n", s.Count)
	fmt.Fprintf(w, "Total ports added:\t%d\n", s.TotalAdded)
	fmt.Fprintf(w, "Total ports removed:\t%d\n", s.TotalRemoved)
	fmt.Fprintf(w, "Scan errors:\t%d\n", s.TotalErrors)
	fmt.Fprintf(w, "Avg scan duration:\t%s\n", roundDuration(s.AvgDuration))
	w.Flush()
}

// PrintLatest writes the most recent scan metrics to the writer.
func (r *MetricsReporter) PrintLatest() {
	sm, ok := r.collector.Latest()
	if !ok {
		fmt.Fprintln(r.out, "no scan data available")
		return
	}
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "=== latest scan ===")
	fmt.Fprintf(w, "Timestamp:\t%s\n", sm.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Duration:\t%s\n", roundDuration(sm.ScanDuration))
	fmt.Fprintf(w, "Ports found:\t%d\n", sm.PortsFound)
	fmt.Fprintf(w, "Ports added:\t%d\n", sm.PortsAdded)
	fmt.Fprintf(w, "Ports removed:\t%d\n", sm.PortsRemoved)
	fmt.Fprintf(w, "Error:\t%v\n", sm.ScanError)
	w.Flush()
}

func roundDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}
