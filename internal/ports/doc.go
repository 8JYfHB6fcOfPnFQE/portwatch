// Package ports provides primitives for scanning, filtering, and tracking
// the state of open network ports on the local host.
//
// # Scanner
//
// Scanner performs a TCP/UDP port sweep over a configurable range and returns
// the set of ports found to be open.
//
// # Filter
//
// Filter allows callers to exclude specific ports or protocols from scan
// results before they are processed further.
//
// # History
//
// History compares successive scan results and surfaces the diff — ports that
// have been opened or closed since the previous scan.
//
// # SnapshotStore
//
// SnapshotStore persists port state to a JSON file so that the daemon can
// resume with a known baseline after a restart, avoiding spurious alerts on
// startup.
package ports
