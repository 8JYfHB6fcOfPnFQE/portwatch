package main

import (
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// applyDiscord wraps next with a DiscordNotifier when Discord is enabled in cfg.
func applyDiscord(cfg config.DiscordConfig, next alert.Notifier) alert.Notifier {
	if !cfg.Enabled {
		return next
	}
	return monitor.NewDiscordNotifier(cfg.WebhookURL, nil, next)
}
