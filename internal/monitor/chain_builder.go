package monitor

import (
	"io"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
)

// ChainBuilder constructs a notifier chain from configuration.
// Each stage wraps the next, so notifications flow through the
// outermost wrapper first and eventually reach the terminal writer.
type ChainBuilder struct {
	cfg    *config.Config
	writer io.Writer
}

// NewChainBuilder returns a builder seeded with the given config and
// terminal output writer.
func NewChainBuilder(cfg *config.Config, w io.Writer) *ChainBuilder {
	return &ChainBuilder{cfg: cfg, writer: w}
}

// Build assembles and returns the head of the notifier chain.
// Stages are added in innermost-to-outermost order so that
// rate-limiting, dedup, and retry wrap the core senders.
func (b *ChainBuilder) Build() alert.Notifier {
	// Terminal: formatted output to writer
	var head alert.Notifier = alert.NewFormattedNotifier(b.writer, b.cfg.Output.Format)

	// Ops instrumentation (innermost wrapper around the real senders)
	if b.cfg.Ops.Enabled {
		head = NewOpsNotifier(head, b.cfg.Ops)
	}

	// Retry on transient failures
	if b.cfg.Retry.Attempts > 0 {
		head = NewRetryNotifier(head, b.cfg.Retry.Attempts, time.Duration(b.cfg.Retry.DelayMs)*time.Millisecond)
	}

	// Deduplication window
	if b.cfg.Dedup.TTL > 0 {
		head = NewDedupNotifier(head, b.cfg.Dedup.TTL)
	}

	// Per-port rate limiting
	if b.cfg.Rate.Limit > 0 {
		head = NewRateNotifier(head, b.cfg.Rate)
	}

	// Silence store suppression
	if b.cfg.SilenceFile != "" {
		store := NewSilenceStore(b.cfg.SilenceFile)
		head = NewSilenceNotifier(head, store)
	}

	// Static label attachment
	if len(b.cfg.Labels) > 0 {
		head = NewLabelNotifier(head, b.cfg.Labels)
	}

	return head
}
