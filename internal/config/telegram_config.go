package config

import "fmt"

// TelegramConfig holds settings for the Telegram notifier.
type TelegramConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	ChatID  string `yaml:"chat_id"`
}

// Validate returns an error if the config is enabled but incomplete.
func (c TelegramConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Token == "" {
		return fmt.Errorf("telegram: token is required")
	}
	if c.ChatID == "" {
		return fmt.Errorf("telegram: chat_id is required")
	}
	return nil
}
