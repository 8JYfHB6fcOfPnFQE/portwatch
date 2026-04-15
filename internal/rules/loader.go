package rules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level structure of a portwatch rules YAML file.
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// LoadFromFile reads and parses a YAML rules file, returning a validated Matcher.
func LoadFromFile(path string) (*Matcher, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening rules file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing rules file %q: %w", path, err)
	}

	if len(cfg.Rules) == 0 {
		return nil, fmt.Errorf("rules file %q contains no rules", path)
	}

	matcher, err := NewMatcher(cfg.Rules)
	if err != nil {
		return nil, fmt.Errorf("validating rules: %w", err)
	}
	return matcher, nil
}
