package main

import (
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// applyEmail wraps next with an email notifier when email is enabled in cfg.
func applyEmail(cfg config.EmailConfig, next alert.Notifier) alert.Notifier {
	if !cfg.Enabled {
		return next
	}
	emailCfg := monitor.EmailConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
		To:       cfg.To,
	}
	return monitor.NewEmailNotifier(emailCfg, next)
}
