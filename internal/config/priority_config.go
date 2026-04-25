package config

import "fmt"

var validPriorities = map[string]bool{
	"low":      true,
	"medium":   true,
	"high":     true,
	"critical": true,
}

// PriorityConfig holds keyword-to-priority mappings used by PriorityNotifier.
type PriorityConfig struct {
	Enabled         bool                `yaml:"enabled"`
	DefaultPriority string              `yaml:"default_priority"`
	Rules           map[string][]string `yaml:"rules"`
}

// DefaultPriorityConfig returns a sensible default configuration.
func DefaultPriorityConfig() PriorityConfig {
	return PriorityConfig{
		Enabled:         false,
		DefaultPriority: "low",
		Rules:           map[string][]string{},
	}
}

// Validate checks that priority levels and keywords are well-formed.
func (c PriorityConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DefaultPriority != "" && !validPriorities[c.DefaultPriority] {
		return fmt.Errorf("priority_config: invalid default_priority %q", c.DefaultPriority)
	}
	for level, keywords := range c.Rules {
		if !validPriorities[level] {
			return fmt.Errorf("priority_config: unknown priority level %q", level)
		}
		for _, kw := range keywords {
			if kw == "" {
				return fmt.Errorf("priority_config: empty keyword under level %q", level)
			}
		}
	}
	return nil
}
