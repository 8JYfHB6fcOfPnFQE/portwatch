package monitor

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func makeStates(portNums ...int) []ports.PortState {
	out := make([]ports.PortState, len(portNums))
	for i, p := range portNums {
		out[i] = ports.PortState{Port: p, Proto: "tcp"}
	}
	return out
}

func TestMetricsHandler_RecordsOnSuccess(t *testing.T) {
	col := ports.NewMetricsCollector(10)
	first := makeStates(80, 443)
	h := NewMetricsHandler(col, func(_ context.Context) ([]ports.PortState, error) {
		return first, nil
	})

	_, err := h.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sm, ok := col.Latest()
	if !ok {
		t.Fatal("expected metrics to be recorded")
	}
	if sm.PortsFound != 2 {
		t.Errorf("expected 2 ports found, got %d", sm.PortsFound)
	}
	if sm.ScanError {
		t.Error("expected no scan error")
	}
}

func TestMetricsHandler_RecordsAddedRemoved(t *testing.T) {
	col := ports.NewMetricsCollector(10)
	calls := 0
	scans := [][]ports.PortState{
		makeStates(80, 443),
		makeStates(80, 8080),
	}
	h := NewMetricsHandler(col, func(_ context.Context) ([]ports.PortState, error) {
		s := scans[calls]
		calls++
		return s, nil
	})

	h.Run(context.Background())
	h.Run(context.Background())

	sm, _ := col.Latest()
	if sm.PortsAdded != 1 {
		t.Errorf("expected 1 added, got %d", sm.PortsAdded)
	}
	if sm.PortsRemoved != 1 {
		t.Errorf("expected 1 removed, got %d", sm.PortsRemoved)
	}
}

func TestMetricsHandler_RecordsError(t *testing.T) {
	col := ports.NewMetricsCollector(10)
	h := NewMetricsHandler(col, func(_ context.Context) ([]ports.PortState, error) {
		return nil, errors.New("scan failed")
	})

	_, err := h.Run(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	sm, _ := col.Latest()
	if !sm.ScanError {
		t.Error("expected ScanError to be true")
	}
}
