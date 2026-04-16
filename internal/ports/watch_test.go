package ports_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

type mockScanner struct {
	results [][]ports.PortState
	call    int
}

func (m *mockScanner) Scan() ([]ports.PortState, error) {
	if m.call >= len(m.results) {
		return m.results[len(m.results)-1], nil
	}
	r := m.results[m.call]
	m.call++
	return r, nil
}

func TestWatcher_EmitsAddedPorts(t *testing.T) {
	ms := &mockScanner{
		results: [][]ports.PortState{
			{},
			{{Port: 8080, Proto: "tcp"}},
		},
	}
	h := ports.NewHistory()
	w := ports.NewWatcher(ms, h, ports.WatchConfig{Interval: 20 * time.Millisecond})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	ch := w.Watch(ctx)
	var got ports.ChangeEvent
	for e := range ch {
		got = e
		break
	}
	if len(got.Added) == 0 {
		t.Fatal("expected added ports")
	}
	if got.Added[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", got.Added[0].Port)
	}
}

func TestWatcher_EmitsRemovedPorts(t *testing.T) {
	ms := &mockScanner{
		results: [][]ports.PortState{
			{{Port: 9090, Proto: "tcp"}},
			{},
		},
	}
	h := ports.NewHistory()
	// prime history
	h.Diff([]ports.PortState{{Port: 9090, Proto: "tcp"}})
	w := ports.NewWatcher(ms, h, ports.WatchConfig{Interval: 20 * time.Millisecond})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	ch := w.Watch(ctx)
	var got ports.ChangeEvent
	for e := range ch {
		got = e
		break
	}
	if len(got.Removed) == 0 {
		t.Fatal("expected removed ports")
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	ms := &mockScanner{results: [][]ports.PortState{{}}}
	h := ports.NewHistory()
	w := ports.NewWatcher(ms, h, ports.WatchConfig{Interval: 10 * time.Millisecond})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch := w.Watch(ctx)
	_, open := <-ch
	if open {
		t.Error("expected channel to be closed")
	}
}
