package monitor

import (
	"testing"
)

func TestNewTagSet_ContainsAddedTags(t *testing.T) {
	ts := NewTagSet([]string{"infra", "dev"})
	if !ts.Contains("infra") {
		t.Error("expected infra to be in set")
	}
	if !ts.Contains("dev") {
		t.Error("expected dev to be in set")
	}
	if ts.Contains("prod") {
		t.Error("prod should not be in set")
	}
}

func TestTagFilter_IsSuppressed_Match(t *testing.T) {
	f := NewTagFilter([]string{"infra", "dev"})
	if !f.IsSuppressed([]string{"web", "infra"}) {
		t.Error("expected suppressed due to infra tag")
	}
}

func TestTagFilter_IsSuppressed_NoMatch(t *testing.T) {
	f := NewTagFilter([]string{"infra"})
	if f.IsSuppressed([]string{"web", "prod"}) {
		t.Error("expected not suppressed")
	}
}

func TestTagFilter_IsSuppressed_EmptyEventTags(t *testing.T) {
	f := NewTagFilter([]string{"infra"})
	if f.IsSuppressed([]string{}) {
		t.Error("empty event tags should not be suppressed")
	}
}

func TestTagFilter_IsSuppressed_EmptySuppressedSet(t *testing.T) {
	f := NewTagFilter([]string{})
	if f.IsSuppressed([]string{"infra", "dev"}) {
		t.Error("no suppressed tags configured, should not suppress")
	}
}
