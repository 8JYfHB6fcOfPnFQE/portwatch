package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestFilterConfig_Validate_Disabled(t *testing.T) {
	fc := config.FilterConfig{Enabled: false, Patterns: []string{""}}
	if err := fc.Validate(); err != nil {
		t.Errorf("disabled config should always pass validation, got: %v", err)
	}
}

func TestFilterConfig_Validate_Valid(t *testing.T) {
	fc := config.FilterConfig{
		Enabled:  true,
		Patterns: []string{"tcp", "0.0.0.0"},
	}
	if err := fc.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}

func TestFilterConfig_Validate_EmptyPattern(t *testing.T) {
	fc := config.FilterConfig{
		Enabled:  true,
		Patterns: []string{"valid", ""},
	}
	if err := fc.Validate(); err == nil {
		t.Error("expected error for empty pattern, got nil")
	}
}

func TestFilterConfig_Default_Disabled(t *testing.T) {
	fc := config.DefaultFilterConfig()
	if fc.Enabled {
		t.Error("default FilterConfig should be disabled")
	}
}

func TestFilterConfig_Default_EmptyPatterns(t *testing.T) {
	fc := config.DefaultFilterConfig()
	if len(fc.Patterns) != 0 {
		t.Errorf("expected empty patterns slice, got %v", fc.Patterns)
	}
}
