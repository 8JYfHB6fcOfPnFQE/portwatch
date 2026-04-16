package monitor_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

type fakeScanner struct{ states []ports.PortState }

func (f *fakeScanner) Scan() ([]ports.PortState, error) { return f.states, nil }

type capturingNotifier struct{ events []alert.Event }

func (c *capturingNotifier) Send(e alert.Event) error {
	c.events = append(c.events, e)
	return nil
}

func TestWatcherHandler_AlertsOnNewPort(t *testing.T) {
	scanner := &fakeScanner{states: []ports.PortState{{Port: 3000, Proto: "tcp"}}}
	h := ports.NewHistory()
	w := ports.NewWatcher(scanner, h, ports.WatchConfig{Interval: 20 * time.Millisecond})
	n := &capturingNotifier{}
	wh := monitor.NewWatcherHandler(w, n, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	wh.Run(ctx)
	if len(n.events) == 0 {
		t.Fatal("expected at least one alert event")
	}
	if n.events[0].Kind != "opened" {
		t.Errorf("expected 'opened', got %q", n.events[0].Kind)
	}
}

func TestWatcherHandler_IgnoreRule(t *testing.T) {
	scanner := &fakeScanner{states: []ports.PortState{{Port: 22, Proto: "tcp"}}}
	h := ports.NewHistory()
	w := ports.NewWatcher(scanner, h, ports.WatchConfig{Interval: 20 * time.Millisecond})
	n := &capturingNotifier{}
	m, _ := rules.NewMatcher([]rules.Rule{{Port: 22, Proto: "tcp", Action: "ignore"}})
	wh := monitor.NewWatcherHandler(w, n, m)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	wh.Run(ctx)
	if len(n.events) != 0 {
		t.Errorf("expected no alerts for ignored port, got %d", len(n.events))
	}
}
