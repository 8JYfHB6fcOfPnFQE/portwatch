package ports

import (
	"testing"
)

func ps(proto string, port int) PortState {
	return PortState{Proto: proto, Port: port}
}

func TestHistory_FirstDiff_AllNew(t *testing.T) {
	h := NewHistory()
	changes := h.Diff([]PortState{ps("tcp", 80), ps("tcp", 443)})
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.OldState != "closed" || c.NewState != "open" {
			t.Errorf("unexpected change states: %+v", c)
		}
	}
}

func TestHistory_DetectsClosedPort(t *testing.T) {
	h := NewHistory()
	h.Diff([]PortState{ps("tcp", 80), ps("tcp", 8080)})

	changes := h.Diff([]PortState{ps("tcp", 80)})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Port.Port != 8080 || changes[0].NewState != "closed" {
		t.Errorf("expected port 8080 closed, got %+v", changes[0])
	}
}

func TestHistory_DetectsOpenedPort(t *testing.T) {
	h := NewHistory()
	h.Diff([]PortState{ps("tcp", 80)})

	changes := h.Diff([]PortState{ps("tcp", 80), ps("tcp", 9090)})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Port.Port != 9090 || changes[0].NewState != "open" {
		t.Errorf("expected port 9090 opened, got %+v", changes[0])
	}
}

func TestHistory_NoChanges(t *testing.T) {
	h := NewHistory()
	h.Diff([]PortState{ps("tcp", 80)})
	changes := h.Diff([]PortState{ps("tcp", 80)})
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestHistory_Snapshot(t *testing.T) {
	h := NewHistory()
	h.Diff([]PortState{ps("tcp", 22), ps("udp", 53)})
	snap := h.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected snapshot length 2, got %d", len(snap))
	}
}

func TestHistory_EmptyDiff(t *testing.T) {
	h := NewHistory()
	changes := h.Diff(nil)
	if len(changes) != 0 {
		t.Errorf("expected no changes on empty diff, got %d", len(changes))
	}
}
