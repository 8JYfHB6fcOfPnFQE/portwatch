package monitor

import (
	"regexp"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// RedactNotifier masks sensitive field values in events before forwarding them
// to the next notifier in the chain. It supports redacting by field name
// (exact match against event metadata keys) and by pattern (regex applied to
// the formatted event address and message).
//
// Redacted values are replaced with the configured placeholder (default: "***").
type RedactNotifier struct {
	next        alert.Notifier
	fields      map[string]struct{}
	patterns    []*regexp.Regexp
	placeholder string
}

// RedactConfig holds the configuration for a RedactNotifier.
type RedactConfig struct {
	// Fields lists metadata key names whose values will be replaced.
	Fields []string
	// Patterns lists regular expressions; any match in Addr or message is replaced.
	Patterns []string
	// Placeholder is the replacement string (defaults to "***").
	Placeholder string
}

// NewRedactNotifier creates a RedactNotifier that masks the given fields and
// pattern matches before forwarding each event to next.
func NewRedactNotifier(cfg RedactConfig, next alert.Notifier) (*RedactNotifier, error) {
	ph := cfg.Placeholder
	if ph == "" {
		ph = "***"
	}

	fields := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		fields[strings.ToLower(f)] = struct{}{}
	}

	var patterns []*regexp.Regexp
	for _, raw := range cfg.Patterns {
		re, err := regexp.Compile(raw)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, re)
	}

	return &RedactNotifier{
		next:        next,
		fields:      fields,
		patterns:    patterns,
		placeholder: ph,
	}, nil
}

// Send redacts the event and forwards it. The original event is never mutated;
// a shallow copy is made before any modifications.
func (r *RedactNotifier) Send(ev alert.Event) error {
	ev = r.redact(ev)
	if r.next != nil {
		return r.next.Send(ev)
	}
	return nil
}

// redact returns a copy of ev with sensitive data replaced.
func (r *RedactNotifier) redact(ev alert.Event) alert.Event {
	// Redact metadata fields by key name.
	if len(r.fields) > 0 && len(ev.Meta) > 0 {
		newMeta := make(map[string]string, len(ev.Meta))
		for k, v := range ev.Meta {
			if _, masked := r.fields[strings.ToLower(k)]; masked {
				newMeta[k] = r.placeholder
			} else {
				newMeta[k] = v
			}
		}
		ev.Meta = newMeta
	}

	// Redact patterns in Addr.
	for _, re := range r.patterns {
		ev.Addr = re.ReplaceAllString(ev.Addr, r.placeholder)
	}

	return ev
}
