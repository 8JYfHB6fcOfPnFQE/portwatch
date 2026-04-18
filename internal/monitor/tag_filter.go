package monitor

// TagFilter suppresses alerts for ports that match a user-defined tag list.
// Tags are arbitrary string labels attached to rules; this allows operators
// to mute whole categories (e.g. "infra", "dev") without editing each rule.

// TagSet is a set of tag strings.
type TagSet map[string]struct{}

// NewTagSet builds a TagSet from a slice of strings.
func NewTagSet(tags []string) TagSet {
	s := make(TagSet, len(tags))
	for _, t := range tags {
		s[t] = struct{}{}
	}
	return s
}

// Contains returns true when the tag is present in the set.
func (ts TagSet) Contains(tag string) bool {
	_, ok := ts[tag]
	return ok
}

// TagFilter decides whether an event should be suppressed based on tags.
type TagFilter struct {
	suppressed TagSet
}

// NewTagFilter creates a TagFilter that will suppress events carrying any of
// the provided tags.
func NewTagFilter(suppressedTags []string) *TagFilter {
	return &TagFilter{suppressed: NewTagSet(suppressedTags)}
}

// IsSuppressed returns true when at least one of the event tags is in the
// suppressed set.
func (f *TagFilter) IsSuppressed(eventTags []string) bool {
	for _, t := range eventTags {
		if f.suppressed.Contains(t) {
			return true
		}
	}
	return false
}
