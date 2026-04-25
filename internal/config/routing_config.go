package config

import (
	"errors"
	"strings"
)

// RoutingConfig holds the configuration for the routing notifier.
type RoutingConfig struct {
	// Enabled controls whether event routing is active.
	Enabled bool `yaml:"enabled"`
	// Rules defines the ordered list of routing rules.
	Rules []RoutingRuleConfig `yaml:"rules"`
}

// RoutingRuleConfig represents a single routing rule in configuration.
type RoutingRuleConfig struct {
	// Field is the event attribute to match against (proto, action, or meta key).
	Field string `yaml:"field"`
	// Value is the substring to match (case-insensitive).
	Value string `yaml:"value"`
	// Destination names the notifier target (e.g. "slack", "email", "webhook").
	Destination string `yaml:"destination"`
}

// DefaultRoutingConfig returns a disabled RoutingConfig with no rules.
func DefaultRoutingConfig() RoutingConfig {
	return RoutingConfig{Enabled: false}
}

// Validate checks that each rule has the required fields set.
func (rc RoutingConfig) Validate() error {
	if !rc.Enabled {
		return nil
	}
	for i, r := range rc.Rules {
		if strings.TrimSpace(r.Field) == "" {
			return errors.New("routing rule missing field at index " + itoa(i))
		}
		if strings.TrimSpace(r.Value) == "" {
			return errors.New("routing rule missing value at index " + itoa(i))
		}
		if strings.TrimSpace(r.Destination) == "" {
			return errors.New("routing rule missing destination at index " + itoa(i))
		}
	}
	return nil
}

// itoa converts an int to its decimal string representation.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
