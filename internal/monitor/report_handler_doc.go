// Package monitor provides the core monitoring loop and supporting handlers
// for portwatch.
//
// ReportHandler
//
// ReportHandler generates an on-demand tabular report of all currently open
// ports detected by the scanner. Each row includes:
//
//   - Protocol (tcp/udp)
//   - Port number
//   - State (e.g. LISTEN)
//   - Process name and PID (via Enricher, if available)
//   - Timestamp of the scan
//
// Usage:
//
//	h := monitor.NewReportHandler(scanner, enricher, os.Stdout)
//	if err := h.Print(); err != nil {
//		log.Fatal(err)
//	}
package monitor
