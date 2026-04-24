package config

import (
	"fmt"
	"time"
)

// RateConfig controls per-window alert rate limiting.
type RateConfig struct {
	// Enabled turns rate limiting on or off.
	Enabled bool `yaml:"enabled"`

	// Limit is the maximum number of alerts forwarded per window.
	// A value of 0 means unlimited.
	Limit int `yaml:"limit"`

	// Window is the duration of each rate-limit window (e.g. "1m", "30s").
	Window time.Duration `yaml:"window"`
}

// DefaultRateConfig returns a sensible default rate configuration.
func DefaultRateConfig() RateConfig {
	return RateConfig{
		Enabled: false,
		Limit:   60,
		Window:  time.Minute,
	}
}

// Validate returns an error if the RateConfig contains invalid values.
func (r RateConfig) Validate() error {
	if !r.Enabled {
		return nil
	}
	if r.Limit < 0 {
		return fmt.Errorf("rate config: limit must be >= 0, got %d", r.Limit)
	}
	if r.Window <= 0 {
		return fmt.Errorf("rate config: window must be positive, got %s", r.Window)
	}
	return nil
}
