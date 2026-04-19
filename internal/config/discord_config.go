package config

import "fmt"

// DiscordConfig holds settings for the Discord webhook notifier.
type DiscordConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

// Validate returns an error if the config is enabled but incomplete.
func (d DiscordConfig) Validate() error {
	if !d.Enabled {
		return nil
	}
	if d.WebhookURL == "" {
		return fmt.Errorf("discord: webhook_url is required when enabled")
	}
	return nil
}
