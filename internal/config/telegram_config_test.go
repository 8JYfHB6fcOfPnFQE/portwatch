package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validTelegramConfig() config.TelegramConfig {
	return config.TelegramConfig{
		Enabled: true,
		Token:   "123456:ABC-DEF",
		ChatID:  "-100123456789",
	}
}

func TestTelegramConfig_Validate_Disabled(t *testing.T) {
	c := config.TelegramConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected nil for disabled config, got %v", err)
	}
}

func TestTelegramConfig_Validate_Valid(t *testing.T) {
	if err := validTelegramConfig().Validate(); err != nil {
		t.Errorf("expected valid config to pass, got %v", err)
	}
}

func TestTelegramConfig_Validate_MissingToken(t *testing.T) {
	c := validTelegramConfig()
	c.Token = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing token")
	}
}

func TestTelegramConfig_Validate_MissingChatID(t *testing.T) {
	c := validTelegramConfig()
	c.ChatID = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing chat_id")
	}
}
