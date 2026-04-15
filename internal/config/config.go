package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"scan_interval"`
	PortRange    PortRange     `yaml:"port_range"`
	Output       OutputConfig  `yaml:"output"`
	RulesFile    string        `yaml:"rules_file"`
}

// PortRange defines the inclusive range of ports to monitor.
type PortRange struct {
	From int `yaml:"from"`
	To   int `yaml:"to"`
}

// OutputConfig controls how alerts are emitted.
type OutputConfig struct {
	Format string `yaml:"format"` // "text" or "json"
	File   string `yaml:"file"`   // empty means stdout
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval: 30 * time.Second,
		PortRange: PortRange{
			From: 1,
			To:   65535,
		},
		Output: OutputConfig{
			Format: "text",
		},
	}
}

// Load reads a YAML config file and merges it over the defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are coherent.
func (c *Config) Validate() error {
	if c.ScanInterval <= 0 {
		return fmt.Errorf("config: scan_interval must be positive")
	}
	if c.PortRange.From < 1 || c.PortRange.From > 65535 {
		return fmt.Errorf("config: port_range.from must be between 1 and 65535")
	}
	if c.PortRange.To < 1 || c.PortRange.To > 65535 {
		return fmt.Errorf("config: port_range.to must be between 1 and 65535")
	}
	if c.PortRange.From > c.PortRange.To {
		return fmt.Errorf("config: port_range.from must be <= port_range.to")
	}
	if c.Output.Format != "text" && c.Output.Format != "json" {
		return fmt.Errorf("config: output.format must be \"text\" or \"json\", got %q", c.Output.Format)
	}
	return nil
}
