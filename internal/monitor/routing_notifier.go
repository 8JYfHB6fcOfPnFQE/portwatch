package monitor

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// RouteRule maps a match condition to a destination notifier.
type RouteRule struct {
	// Field is the event field to inspect: "proto", "action", or a meta key.
	Field string
	// Value is the case-insensitive substring to match against the field.
	Value string
	// Notifier receives events that satisfy this rule.
	Notifier alert.Notifier
}

// RoutingNotifier dispatches each event to the first matching RouteRule.
// If no rule matches, the event is forwarded to the fallback notifier (if set).
type RoutingNotifier struct {
	rules    []RouteRule
	fallback alert.Notifier
}

// NewRoutingNotifier creates a RoutingNotifier with the given rules and an
// optional fallback notifier for unmatched events.
func NewRoutingNotifier(rules []RouteRule, fallback alert.Notifier) *RoutingNotifier {
	return &RoutingNotifier{rules: rules, fallback: fallback}
}

// Send evaluates each RouteRule in order and dispatches the event to the first
// match. If no rule matches and a fallback is configured, the event is sent
// there instead.
func (r *RoutingNotifier) Send(e alert.Event) error {
	for _, rule := range r.rules {
		if r.matches(e, rule) {
			if rule.Notifier != nil {
				return rule.Notifier.Send(e)
			}
			return nil
		}
	}
	if r.fallback != nil {
		return r.fallback.Send(e)
	}
	return nil
}

func (r *RoutingNotifier) matches(e alert.Event, rule RouteRule) bool {
	var candidate string
	switch strings.ToLower(rule.Field) {
	case "proto":
		candidate = e.Proto
	case "action":
		candidate = e.Action
	default:
		// treat as a metadata key lookup
		if e.Meta != nil {
			candidate = e.Meta[rule.Field]
		}
	}
	return strings.Contains(strings.ToLower(candidate), strings.ToLower(rule.Value))
}
