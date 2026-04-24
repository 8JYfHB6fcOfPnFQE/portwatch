package config

import "fmt"

// OpsConfig controls the operational instrumentation notifier.
type OpsConfig struct {
	// Enabled turns the ops instrumentation layer on or off.
	Enabled bool `yaml:"enabled"`

	// LogFailed controls whether failed sends are written to the logger.
	LogFailed bool `yaml:"log_failed"`

	// LatencyWarnMs triggers a warning log when a single send exceeds this
	// threshold (0 = disabled).
	LatencyWarnMs int64 `yaml:"latency_warn_ms"`
}

// Validate returns an error when OpsConfig contains invalid values.
func (o OpsConfig) Validate() error {
	if o.LatencyWarnMs < 0 {
		return fmt.Errorf("ops: latency_warn_ms must be >= 0, got %d", o.LatencyWarnMs)
	}
	return nil
}

// DefaultOpsConfig returns a safe default configuration.
func DefaultOpsConfig() OpsConfig {
	return OpsConfig{
		Enabled:       true,
		LogFailed:     true,
		LatencyWarnMs: 0,
	}
}
