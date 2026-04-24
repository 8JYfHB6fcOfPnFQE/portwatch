package config

import "fmt"

// FilterConfig holds configuration for the filter notifier which drops events
// whose content matches any of the supplied substring patterns.
type FilterConfig struct {
	// Enabled controls whether the filter notifier is active.
	Enabled bool `yaml:"enabled"`

	// Patterns is a list of case-insensitive substrings. An event is dropped
	// when its string representation or any meta value contains a match.
	Patterns []string `yaml:"patterns"`
}

// DefaultFilterConfig returns a FilterConfig with filtering disabled.
func DefaultFilterConfig() FilterConfig {
	return FilterConfig{
		Enabled:  false,
		Patterns: []string{},
	}
}

// Validate returns an error if the configuration is inconsistent.
func (fc FilterConfig) Validate() error {
	if !fc.Enabled {
		return nil
	}
	for i, p := range fc.Patterns {
		if p == "" {
			return fmt.Errorf("filter: pattern at index %d must not be empty", i)
		}
	}
	return nil
}
