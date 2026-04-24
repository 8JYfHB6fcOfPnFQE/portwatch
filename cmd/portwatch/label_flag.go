package main

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

// applyLabels wraps next with a LabelNotifier when one or more --label flags
// have been provided. Labels are expected in "key=value" format; malformed
// entries are silently skipped.
//
// Example CLI usage:
//
//	portwatch --label env=prod --label region=us-east-1
func applyLabels(rawLabels []string, next alert.Notifier) alert.Notifier {
	if len(rawLabels) == 0 {
		return next
	}

	labels := make(map[string]string, len(rawLabels))
	for _, raw := range rawLabels {
		k, v, ok := strings.Cut(raw, "=")
		if !ok || strings.TrimSpace(k) == "" {
			continue
		}
		labels[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}

	if len(labels) == 0 {
		return next
	}
	return monitor.NewLabelNotifier(labels, next)
}
