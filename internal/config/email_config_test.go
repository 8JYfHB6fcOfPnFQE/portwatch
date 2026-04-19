package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validEmailConfig() config.EmailConfig {
	return config.EmailConfig{
		Enabled:  true,
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "pass",
		From:     "alert@example.com",
		To:       []string{"ops@example.com"},
	}
}

func TestEmailConfig_Validate_Disabled(t *testing.T) {
	cfg := config.EmailConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for disabled config, got %v", err)
	}
}

func TestEmailConfig_Validate_Valid(t *testing.T) {
	if err := validEmailConfig().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmailConfig_Validate_MissingHost(t *testing.T) {
	cfg := validEmailConfig()
	cfg.Host = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestEmailConfig_Validate_InvalidPort(t *testing.T) {
	cfg := validEmailConfig()
	cfg.Port = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestEmailConfig_Validate_MissingFrom(t *testing.T) {
	cfg := validEmailConfig()
	cfg.From = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestEmailConfig_Validate_NoRecipients(t *testing.T) {
	cfg := validEmailConfig()
	cfg.To = nil
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for no recipients")
	}
}
