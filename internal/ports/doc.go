// Package ports provides primitives for discovering, tracking, and filtering
// open network ports on the local host.
//
// # Scanner
//
// Scanner performs active TCP/UDP port scans over a configurable range using
// net.Dial, suitable for cross-platform use.
//
// # ProcReader
//
// ProcReader reads listening ports directly from /proc/net/tcp and
// /proc/net/tcp6 on Linux systems. It is faster than active scanning and
// does not require elevated privileges beyond read access to /proc.
//
// # Filter
//
// Filter applies exclusion rules so that well-known or explicitly ignored
// ports and protocols are suppressed before alerting.
//
// # History
//
// History compares successive port snapshots and emits the diff — newly
// opened ports and recently closed ports — for downstream alerting.
//
// # SnapshotStore
//
// SnapshotStore persists port snapshots to disk so that portwatch can
// detect changes across restarts without a false-positive flood on startup.
package ports
