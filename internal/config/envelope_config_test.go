package config_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestEnvelopeConfig_Validate_Disabled(t *testing.T) {
	cfg := config.EnvelopeConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for disabled config, got %v", err)
	}
}

func TestEnvelopeConfig_Validate_ValidSource(t *testing.T) {
	cfg := config.EnvelopeConfig{Enabled: true, Source: "prod-node-1"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnvelopeConfig_Validate_EmptySource_OK(t *testing.T) {
	cfg := config.EnvelopeConfig{Enabled: true, Source: ""}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("empty source should be valid, got %v", err)
	}
}

func TestEnvelopeConfig_Validate_LongSource_Error(t *testing.T) {
	cfg := config.EnvelopeConfig{
		Enabled: true,
		Source:  strings.Repeat("x", 129),
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for source > 128 chars")
	}
}

func TestEnvelopeConfig_Default_Disabled(t *testing.T) {
	cfg := config.DefaultEnvelopeConfig()
	if cfg.Enabled {
		t.Fatal("default envelope config should be disabled")
	}
}

func TestEnvelopeConfig_Default_EmptySource(t *testing.T) {
	cfg := config.DefaultEnvelopeConfig()
	if cfg.Source != "" {
		t.Fatalf("expected empty source, got %q", cfg.Source)
	}
}
