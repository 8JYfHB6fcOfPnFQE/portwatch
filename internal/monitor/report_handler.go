package monitor

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// ReportHandler prints a formatted port report on demand.
type ReportHandler struct {
	scanner ports.Scanner
	enricher ports.Enricher
	out      io.Writer
}

// NewReportHandler creates a ReportHandler that writes to w (defaults to os.Stdout).
func NewReportHandler(scanner ports.Scanner, enricher ports.Enricher, w io.Writer) *ReportHandler {
	if w == nil {
		w = os.Stdout
	}
	return &ReportHandler{scanner: scanner, enricher: enricher, out: w}
}

// Print scans current ports and writes a tabular report.
func (h *ReportHandler) Print() error {
	states, err := h.scanner.Scan()
	if err != nil {
		return fmt.Errorf("report scan: %w", err)
	}

	tw := tabwriter.NewWriter(h.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROTO\tPORT\tSTATE\tPROCESS\tSCANNED")

	now := time.Now().Format(time.RFC3339)
	for _, s := range states {
		proc := "unknown"
		if h.enricher != nil {
			if info, err := h.enricher.Lookup(s.Inode); err == nil {
				proc = info.String()
			}
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\n", s.Proto, s.Port, s.State, proc, now)
	}
	return tw.Flush()
}
