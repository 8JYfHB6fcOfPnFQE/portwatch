package main

import (
	"net/http"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// applyTelegram wraps next with a TelegramNotifier if Telegram is enabled in cfg.
func applyTelegram(cfg config.TelegramConfig, next alert.Notifier) alert.Notifier {
	if !cfg.Enabled {
		return next
	}
	return monitor.NewTelegramNotifier(cfg.Token, cfg.ChatID, http.DefaultClient, next)
}
