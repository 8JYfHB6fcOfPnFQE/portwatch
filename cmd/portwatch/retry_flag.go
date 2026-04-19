package main

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

// retryConfig holds CLI-level retry settings parsed from flags or config.
type retryConfig struct {
	Enabled  bool
	MaxTries int
	Delay    time.Duration
}

// applyRetry wraps the supplied Notifier with a RetryNotifier when retry is
// enabled.  If disabled the original notifier is returned unchanged.
func applyRetry(n alert.Notifier, cfg retryConfig, logger *log.Logger) alert.Notifier {
	if !cfg.Enabled || cfg.MaxTries <= 1 {
		return n
	}
	delay := cfg.Delay
	if delay <= 0 {
		delay = 500 * time.Millisecond
	}
	return monitor.NewRetryNotifier(n, cfg.MaxTries, delay, logger)
}
