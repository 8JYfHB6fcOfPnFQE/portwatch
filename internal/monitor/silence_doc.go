// Package monitor provides the core monitoring loop and associated handlers
// for portwatch.
//
// # Silence
//
// The silence subsystem allows operators to temporarily suppress alerts for
// specific port/protocol pairs. This is useful during planned maintenance
// windows where a service is intentionally stopped or restarted.
//
// Components:
//
//   - SilenceStore   – thread-safe in-memory store of active silence rules.
//   - SilenceNotifier – alert.Notifier wrapper that drops events for silenced ports.
//   - SilenceHandler  – CLI helper for add / list / purge operations.
//
// Silences are keyed by (port, proto) and expire after a caller-supplied
// duration. Call Purge periodically (e.g. each monitor tick) to reclaim
// memory from expired entries.
package monitor
