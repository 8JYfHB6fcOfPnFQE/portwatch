package config

import "fmt"

// EnvelopeConfig controls the EnvelopeNotifier middleware.
type EnvelopeConfig struct {
	// Enabled turns the envelope notifier on or off.
	Enabled bool `yaml:"enabled"`

	// Source is an arbitrary label embedded in every envelope (e.g. hostname).
	// When empty the envelope.source metadata key is omitted.
	Source string `yaml:"source"`
}

// DefaultEnvelopeConfig returns a sensible default: disabled with no source.
func DefaultEnvelopeConfig() EnvelopeConfig {
	return EnvelopeConfig{
		Enabled: false,
		Source:  "",
	}
}

// Validate checks that the EnvelopeConfig is self-consistent.
func (e EnvelopeConfig) Validate() error {
	if !e.Enabled {
		return nil
	}
	if len(e.Source) > 128 {
		return fmt.Errorf("envelope source label exceeds 128 characters")
	}
	return nil
}
