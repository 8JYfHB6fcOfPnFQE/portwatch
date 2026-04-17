package monitor

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makePorts(n int) []ports.PortState {
	out := make([]ports.PortState, n)
	for i := range out {
		out[i] = ports.PortState{Port: 8000 + i, Proto: "tcp"}
	}
	return out
}

func TestSuppressHandler_WithinWindow_FiltersAll(t *testing.T) {
	h := NewSuppressHandler(10 * time.Second)
	diff := makePorts(3)
	result := h.FilterDiff(diff)
	if len(result) != 0 {
		t.Fatalf("expected 0 events, got %d", len(result))
	}
	if h.Suppressed() != 3 {
		t.Fatalf("expected 3 suppressed, got %d", h.Suppressed())
	}
}

func TestSuppressHandler_AfterWindow_PassesThrough(t *testing.T) {
	h := NewSuppressHandler(0)
	time.Sleep(2 * time.Millisecond)
	diff := makePorts(2)
	result := h.FilterDiff(diff)
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
	if h.Suppressed() != 0 {
		t.Fatalf("expected 0 suppressed, got %d", h.Suppressed())
	}
}

func TestSuppressHandler_IsSuppressed_True(t *testing.T) {
	h := NewSuppressHandler(10 * time.Second)
	if !h.IsSuppressed() {
		t.Fatal("expected handler to be suppressed")
	}
}

func TestSuppressHandler_IsSuppressed_False(t *testing.T) {
	h := NewSuppressHandler(0)
	time.Sleep(2 * time.Millisecond)
	if h.IsSuppressed() {
		t.Fatal("expected handler to not be suppressed")
	}
}

func TestSuppressHandler_Reset_RestartsWindow(t *testing.T) {
	h := NewSuppressHandler(0)
	time.Sleep(2 * time.Millisecond)
	h.Reset()
	if !h.IsSuppressed() {
		t.Fatal("expected suppression after reset")
	}
	if h.Suppressed() != 0 {
		t.Fatalf("expected suppressed count reset to 0")
	}
}

func TestSuppressHandler_AccumulatesSuppressedCount(t *testing.T) {
	h := NewSuppressHandler(10 * time.Second)
	h.FilterDiff(makePorts(2))
	h.FilterDiff(makePorts(3))
	if h.Suppressed() != 5 {
		t.Fatalf("expected 5 suppressed, got %d", h.Suppressed())
	}
}
