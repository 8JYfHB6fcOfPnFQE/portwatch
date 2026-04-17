// Package monitor provides the core monitoring loop and supporting handlers
// for portwatch.
//
// SuppressHandler
//
// SuppressHandler implements a startup quiet window that prevents alert
// storms when the daemon restarts and observes the existing set of open
// ports as "new" changes.
//
// Usage:
//
//	sh := monitor.NewSuppressHandler(30 * time.Second)
//
//	// Inside the monitor tick:
//	diff = sh.FilterDiff(diff)
//	if len(diff) == 0 {
//		return // suppressed
//	}
//
// The suppression window begins at construction time (or after Reset).
// Events dropped during the window are counted and available via Suppressed().
package monitor
