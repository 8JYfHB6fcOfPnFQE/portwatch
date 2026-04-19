package config

import "testing"

func validDiscordConfig() DiscordConfig {
	return DiscordConfig{
		Enabled:    true,
		WebhookURL: "https://discord.com/api/webhooks/123/abc",
	}
}

func TestDiscordConfig_Validate_Disabled(t *testing.T) {
	c := DiscordConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestDiscordConfig_Validate_Valid(t *testing.T) {
	if err := validDiscordConfig().Validate(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestDiscordConfig_Validate_MissingWebhookURL(t *testing.T) {
	c := DiscordConfig{Enabled: true, WebhookURL: ""}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing webhook_url")
	}
}
