// Package monitor provides the DedupNotifier, which wraps any alert.Notifier
// and suppresses duplicate events for the same port, protocol, and event kind
// within a configurable TTL window.
//
// This is useful when a port flaps rapidly or when the monitor loop fires
// faster than human-readable alert cadence requires.
//
// Usage:
//
//	base := alert.NewNotifier(os.Stdout)
//	dedup := monitor.NewDedupNotifier(base, 30*time.Second)
//	// pass dedup wherever alert.Notifier is accepted
//
// Flush() can be called to clear all cached keys, for example on SIGHUP.
package monitor
