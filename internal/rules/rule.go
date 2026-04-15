package rules

import "fmt"

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAlert  Action = "alert"
	ActionIgnore Action = "ignore"
)

// Rule defines a condition to match against an observed port event.
type Rule struct {
	Name    string `yaml:"name"`
	Port    int    `yaml:"port"`
	Proto   string `yaml:"proto"`
	Action  Action `yaml:"action"`
	Comment string `yaml:"comment,omitempty"`
}

// Validate checks that the rule fields are sensible.
func (r *Rule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("rule must have a name")
	}
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("rule %q: port %d out of valid range (1-65535)", r.Name, r.Port)
	}
	if r.Proto != "tcp" && r.Proto != "udp" {
		return fmt.Errorf("rule %q: proto must be 'tcp' or 'udp', got %q", r.Name, r.Proto)
	}
	if r.Action != ActionAlert && r.Action != ActionIgnore {
		return fmt.Errorf("rule %q: action must be 'alert' or 'ignore', got %q", r.Name, r.Action)
	}
	return nil
}

// Matcher evaluates a set of rules against port/proto pairs.
type Matcher struct {
	rules []Rule
}

// NewMatcher creates a Matcher after validating all rules.
func NewMatcher(rules []Rule) (*Matcher, error) {
	for i := range rules {
		if err := rules[i].Validate(); err != nil {
			return nil, err
		}
	}
	return &Matcher{rules: rules}, nil
}

// Match returns the first rule that matches the given port and proto.
// If no rule matches, nil is returned.
func (m *Matcher) Match(port int, proto string) *Rule {
	for i := range m.rules {
		r := &m.rules[i]
		if r.Port == port && r.Proto == proto {
			return r
		}
	}
	return nil
}
