// Package ports provides port scanning, filtering, and watching utilities.
package ports

import (
	"context"
	"time"
)

// WatchConfig holds configuration for the port watcher.
type WatchConfig struct {
	Interval time.Duration
	Filter   *Filter
}

// ChangeEvent represents a detected port change.
type ChangeEvent struct {
	Added   []PortStatePortState
}

// Watcher continuously scans ports and emits change events.
type Watcher struct {
	scanner Scanner
	histfg     WatchConfig
}

// Scanner is the interface for scanning open ports.
type Scanner interface {
	Scan() ([]PortState, error)
}

// New Watcher.
func NewWatcher(scanner Scanner, history *History, cfg WatchConfig) *Watcher {
	return &Watcher{scanner: scanner, history: history, cfg: cfg}
}

// Watch starts watching for port changes, sending events to the returned channel.
// It stops when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan ChangeEvent {
	ch := make(chan ChangeEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ports, err := w.scanner.Scan()
		t		c.Filter != nil {
	var filtered []PortState
					for _, p := range ports {
						if w.cfg.Filter.Allow(p) {
							filtered = append(filtered, p)
						}
					}
					ports = filtered
				}
				added, removed := w.history.Diff(ports)
				if len(added) > 0 || len(removed) > 0 {
					select {
					case ch <- ChangeEvent{Added: added, Removed: removed}:
					default:
					}
				}
			}
		}
	}()
	return ch
}
