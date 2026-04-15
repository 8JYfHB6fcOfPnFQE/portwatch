// Package config loads and validates portwatch configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	Interval      time.Duration `yaml:"interval"`
	PortRange     string        `yaml:"port_range"`
	RulesFile     string        `yaml:"rules_file"`
	OutputFormat  string        `yaml:"output_format"`
	ExcludePorts  []int         `yaml:"exclude_ports"`
	ExcludeProtos []string      `yaml:"exclude_protos"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Interval:      15 * time.Second,
		PortRange:     "1-65535",
		OutputFormat:  "text",
		ExcludePorts:  []int{},
		ExcludeProtos: []string{},
	}
}

// Load reads a YAML config file and merges it with defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks that the Config fields are within acceptable ranges.
func (c *Config) Validate() error {
	if c.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1s, got %s", c.Interval)
	}
	if c.PortRange == "" {
		return errors.New("port_range must not be empty")
	}
	switch c.OutputFormat {
	case "text", "json":
		// valid
	case "":
		c.OutputFormat = "text"
	default:
		return fmt.Errorf("unknown output_format %q: must be \"text\" or \"json\"", c.OutputFormat)
	}
	return nil
}
