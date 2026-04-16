package monitor

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/ports"
)

// BaselineHandler provides CLI-level operations for managing port baselines.
type BaselineHandler struct {
	store  *ports.BaselineStore
	scanner interface {
		Scan() ([]ports.PortState, error)
	}
	out io.Writer
}

// NewBaselineHandler creates a BaselineHandler.
func NewBaselineHandler(store *ports.BaselineStore, scanner interface {
	Scan() ([]ports.PortState, error)
}, out io.Writer) *BaselineHandler {
	if out == nil {
		out = os.Stdout
	}
	return &BaselineHandler{store: store, scanner: scanner, out: out}
}

// Capture scans current ports and saves them as the new baseline.
func (h *BaselineHandler) Capture() error {
	states, err := h.scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}
	if err := h.store.Save(states); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	fmt.Fprintf(h.out, "baseline captured: %d ports saved to %s\n", len(states), h.store.Path())
	return nil
}

// Show prints the current baseline to the configured writer.
func (h *BaselineHandler) Show() error {
	bl, err := h.store.Load()
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	if bl == nil {
		fmt.Fprintln(h.out, "no baseline found")
		return nil
	}
	fmt.Fprintf(h.out, "baseline captured at %s (%d ports):\n", bl.CreatedAt.Format("2006-01-02 15:04:05"), len(bl.Ports))
	for _, ps := range bl.Ports {
		fmt.Fprintf(h.out, "  %s\n", ps)
	}
	return nil
}

// Check compares the current scan against the saved baseline and reports differences.
func (h *BaselineHandler) Check() error {
	bl, err := h.store.Load()
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	if bl == nil {
		fmt.Fprintln(h.out, "no baseline found; run 'portwatch baseline capture' first")
		return nil
	}
	current, err := h.scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}
	added, removed := bl.Diff(current)
	if len(added) == 0 && len(removed) == 0 {
		fmt.Fprintln(h.out, "no changes from baseline")
		return nil
	}
	for _, ps := range added {
		fmt.Fprintf(h.out, "ADDED   %s\n", ps)
	}
	for _, ps := range removed {
		fmt.Fprintf(h.out, "REMOVED %s\n", ps)
	}
	return nil
}
