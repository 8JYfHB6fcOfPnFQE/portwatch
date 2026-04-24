package config

import "testing"

func TestOpsConfig_Validate_Valid(t *testing.T) {
	cfg := DefaultOpsConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpsConfig_Validate_ZeroLatency(t *testing.T) {
	cfg := OpsConfig{Enabled: true, LatencyWarnMs: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("zero latency_warn_ms should be valid: %v", err)
	}
}

func TestOpsConfig_Validate_NegativeLatency(t *testing.T) {
	cfg := OpsConfig{Enabled: true, LatencyWarnMs: -1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative latency_warn_ms")
	}
}

func TestOpsConfig_Default_Enabled(t *testing.T) {
	cfg := DefaultOpsConfig()
	if !cfg.Enabled {
		t.Fatal("default ops config should be enabled")
	}
}

func TestOpsConfig_Default_LogFailed(t *testing.T) {
	cfg := DefaultOpsConfig()
	if !cfg.LogFailed {
		t.Fatal("default ops config should log failed sends")
	}
}

func TestOpsConfig_Disabled_Validate_Passes(t *testing.T) {
	cfg := OpsConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled ops config should still validate: %v", err)
	}
}
