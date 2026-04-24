package main

import (
	"log"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// applyOps wraps notifier with the operational instrumentation layer when the
// feature is enabled in cfg. The shared OpsMetrics pointer is returned so the
// caller can expose it via the metrics reporter.
func applyOps(notifier alert.Notifier, cfg config.OpsConfig, logger *log.Logger) (alert.Notifier, *monitor.OpsMetrics) {
	if !cfg.Enabled {
		return notifier, nil
	}
	m := &monitor.OpsMetrics{}
	on := monitor.NewOpsNotifier(notifier, m, logger)
	return on, m
}
