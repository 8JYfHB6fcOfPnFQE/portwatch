package main

import (
	"os"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// applyEnvelope wraps notifier with an EnvelopeNotifier when the envelope
// feature is enabled in cfg. The source label defaults to the system hostname
// when the config source is empty and the feature is enabled.
func applyEnvelope(cfg config.EnvelopeConfig, notifier alert.Notifier) alert.Notifier {
	if !cfg.Enabled {
		return notifier
	}

	source := cfg.Source
	if source == "" {
		if h, err := os.Hostname(); err == nil {
			source = h
		}
	}

	return monitor.NewEnvelopeNotifier(source, notifier)
}
