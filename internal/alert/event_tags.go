package alert

// Tags attaches an arbitrary set of string labels to an Event.
// Tags are used by downstream filters (e.g. TagNotifier) to suppress or route
// alerts without modifying rule definitions.

// WithTags returns a copy of ev with the given tags appended.
func WithTags(ev Event, tags ...string) Event {
	copy := ev
	copy.Tags = append(append([]string{}, ev.Tags...), tags...)
	return copy
}

// HasTag reports whether ev carries the given tag.
func HasTag(ev Event, tag string) bool {
	for _, t := range ev.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
