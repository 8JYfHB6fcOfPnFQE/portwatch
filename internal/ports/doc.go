// Package ports provides primitives for scanning, tracking, and enriching
// open network port information on the local host.
//
// # Scanner
//
// NewScanner returns a Scanner that lists currently open TCP/UDP ports within
// a configurable port range.
//
// # History
//
// NewHistory tracks successive scans and exposes Diff to surface ports that
// have opened or closed between two snapshots.
//
// # Snapshot
//
// NewSnapshotStore persists port state to disk so that portwatch can survive
// restarts without generating spurious alerts for pre-existing ports.
//
// # Baseline
//
// NewBaselineStore saves a named reference snapshot that operators capture
// intentionally. Baseline.Diff compares any subsequent scan against that
// reference to highlight unexpected additions or removals.
//
// # Filter
//
// NewFilter applies exclusion rules (by port number or protocol) before
// results are forwarded to the alerting layer.
//
// # Enricher
//
// NewEnricher correlates open ports with /proc data to attach process name
// and PID metadata to each PortState.
//
// # RateLimiter
//
// NewRateLimiter suppresses repeated alerts for the same port within a
// configurable cooldown window.
