package config

import "fmt"

// SamplingConfig controls statistical event sampling in the notifier chain.
type SamplingConfig struct {
	// Enabled toggles sampling. When false the notifier is a no-op pass-through.
	Enabled bool `yaml:"enabled"`

	// Rate is the fraction of events to forward (0.0 = drop all, 1.0 = keep all).
	Rate float64 `yaml:"rate"`
}

// DefaultSamplingConfig returns a SamplingConfig with sampling disabled.
func DefaultSamplingConfig() SamplingConfig {
	return SamplingConfig{
		Enabled: false,
		Rate:    1.0,
	}
}

// Validate returns an error if the configuration is invalid.
func (s SamplingConfig) Validate() error {
	if !s.Enabled {
		return nil
	}
	if s.Rate < 0 || s.Rate > 1 {
		return fmt.Errorf("sampling.rate must be between 0.0 and 1.0, got %g", s.Rate)
	}
	return nil
}
