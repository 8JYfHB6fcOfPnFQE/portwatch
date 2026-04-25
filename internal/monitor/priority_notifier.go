package monitor

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Priority levels for events.
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityCritical = "critical"
)

// PriorityNotifier assigns a priority level to each event based on
// configurable keyword patterns, then forwards to the next notifier.
// Events that match a higher-priority pattern are tagged accordingly.
type PriorityNotifier struct {
	next     alert.Notifier
	rules    []priorityRule
	default_ string
}

type priorityRule struct {
	priority string
	keywords []string
}

// NewPriorityNotifier constructs a PriorityNotifier. rules maps a priority
// level (e.g. "high") to a slice of keywords matched case-insensitively
// against the event message. The first matching rule wins.
func NewPriorityNotifier(next alert.Notifier, rules map[string][]string, defaultPriority string) *PriorityNotifier {
	if defaultPriority == "" {
		defaultPriority = PriorityLow
	}
	// Ordered evaluation: critical > high > medium > low
	order := []string{PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow}
	var ordered []priorityRule
	for _, p := range order {
		if kw, ok := rules[p]; ok && len(kw) > 0 {
			ordered = append(ordered, priorityRule{priority: p, keywords: kw})
		}
	}
	return &PriorityNotifier{next: next, rules: ordered, default_: defaultPriority}
}

// Send assigns a priority to the event and forwards it.
func (n *PriorityNotifier) Send(ev alert.Event) error {
	priority := n.resolve(ev)
	ev = alert.WithMeta(ev, "priority", priority)
	if n.next != nil {
		return n.next.Send(ev)
	}
	return nil
}

func (n *PriorityNotifier) resolve(ev alert.Event) string {
	msg := strings.ToLower(ev.Message)
	for _, rule := range n.rules {
		for _, kw := range rule.keywords {
			if strings.Contains(msg, strings.ToLower(kw)) {
				return rule.priority
			}
		}
	}
	return n.default_
}
