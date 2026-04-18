package monitor

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// SilenceHandler exposes CLI sub-commands for managing silences.
type SilenceHandler struct {
	store *SilenceStore
	out   io.Writer
}

// NewSilenceHandler returns a SilenceHandler writing output to w.
func NewSilenceHandler(store *SilenceStore, w io.Writer) *SilenceHandler {
	return &SilenceHandler{store: store, out: w}
}

// Add silences port/proto for dur (e.g. "80", "tcp", 30*time.Minute).
func (h *SilenceHandler) Add(port int, proto string, dur time.Duration) {
	h.store.Add(port, proto, dur)
	fmt.Fprintf(h.out, "silenced %s/%d for %s\n", proto, port, dur)
}

// List prints all active silences.
func (h *SilenceHandler) List() {
	rules := h.store.List()
	if len(rules) == 0 {
		fmt.Fprintln(h.out, "no active silences")
		return
	}
	var sb strings.Builder
	for _, r := range rules {
		sb.WriteString(fmt.Sprintf("%s/%s until %s\n",
			strconv.Itoa(r.Port), r.Proto, r.Deadline.Format(time.RFC3339)))
	}
	fmt.Fprint(h.out, sb.String())
}

// Purge removes expired silences and reports the count removed.
func (h *SilenceHandler) Purge() {
	before := len(h.store.List())
	h.store.Purge()
	after := len(h.store.List())
	fmt.Fprintf(h.out, "purged %d expired silence(s)\n", before-after)
}
