package monitor

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SnapshotHandler saves port state snapshots periodically and on demand.
type SnapshotHandler struct {
	store   *ports.SnapshotStore
	scanner *ports.Scanner
	out     io.Writer
}

// NewSnapshotHandler creates a SnapshotHandler backed by the given store and scanner.
func NewSnapshotHandler(store *ports.SnapshotStore, scanner *ports.Scanner) *SnapshotHandler {
	return &SnapshotHandler{
		store:   store,
		scanner: scanner,
		out:     os.Stdout,
	}
}

// Capture scans current ports and writes a snapshot to disk.
func (h *SnapshotHandler) Capture(ctx context.Context) error {
	states, err := h.scanner.Scan(ctx)
	if err != nil {
		return fmt.Errorf("snapshot scan: %w", err)
	}
	if err := h.store.Save(states); err != nil {
		return fmt.Errorf("snapshot save: %w", err)
	}
	return nil
}

// Show prints the most recently saved snapshot to the handler's writer.
func (h *SnapshotHandler) Show() error {
	states, err := h.store.Load()
	if err != nil {
		return fmt.Errorf("snapshot load: %w", err)
	}
	if len(states) == 0 {
		fmt.Fprintln(h.out, "no snapshot available")
		return nil
	}
	tw := tabwriter.NewWriter(h.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROTO\tPORT\tADDR")
	for _, s := range states {
		fmt.Fprintf(tw, "%s\t%d\t%s\n", s.Proto, s.Port, s.Addr)
	}
	return tw.Flush()
}

// StartPeriodicCapture captures a snapshot every interval until ctx is cancelled.
func (h *SnapshotHandler) StartPeriodicCapture(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = h.Capture(ctx)
		case <-ctx.Done():
			return
		}
	}
}
