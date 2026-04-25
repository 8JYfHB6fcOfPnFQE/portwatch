package config

import "fmt"

// TransformConfig controls the built-in event transformation pipeline.
// Transforms are applied in the order they are listed here before the
// event reaches any downstream notifier.
type TransformConfig struct {
	// UpperCaseAction converts the event action string to upper-case.
	UpperCaseAction bool `yaml:"upper_case_action"`

	// RedactAddr strips the remote address from every event.
	RedactAddr bool `yaml:"redact_addr"`

	// StaticMeta is a map of key/value pairs injected into every event's
	// metadata. Values from this map do NOT override labels set by
	// LabelNotifier; they are applied first in the chain.
	StaticMeta map[string]string `yaml:"static_meta"`
}

// DefaultTransformConfig returns a TransformConfig with safe defaults
// (all transforms disabled, no static metadata).
func DefaultTransformConfig() TransformConfig {
	return TransformConfig{
		StaticMeta: make(map[string]string),
	}
}

// Validate returns an error if any field contains an invalid value.
// Currently the only restriction is that static metadata keys must not
// be empty strings.
func (c TransformConfig) Validate() error {
	for k := range c.StaticMeta {
		if k == "" {
			return fmt.Errorf("transform: static_meta contains an empty key")
		}
	}
	return nil
}

// IsNoOp reports whether the config would produce any mutations.
func (c TransformConfig) IsNoOp() bool {
	return !c.UpperCaseAction && !c.RedactAddr && len(c.StaticMeta) == 0
}
