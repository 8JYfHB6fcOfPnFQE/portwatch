package monitor

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

// WatcherHandler wires a Watcher to an alert Notifier with rule matching.
type WatcherHandler struct {
	watcher  *ports.Watcher
	notifier alert.Notifier
	matcher  *rules.Matcher
}

// NewWatcherHandler creates a WatcherHandler.
func NewWatcherHandler(w *ports.Watcher, n alert.Notifier, m *rules.Matcher) *WatcherHandler {
	return &WatcherHandler{watcher: w, notifier: n, matcher: m}
}

// Run starts processing change events until ctx is cancelled.
func (h *WatcherHandler) Run(ctx context.Context) {
	ch := h.watcher.Watch(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			h.handleEvent(ev)
		}
	}
}

func (h *WatcherHandler) handleEvent(ev ports.ChangeEvent) {
	for _, p := range ev.Added {
		action := h.matchAction(p)
		if action == "ignore" {
			continue
		}
		e := alert.NewEvent("opened", p)
		if err := h.notifier.Send(e); err != nil {
			log.Printf("alert send error: %v", err)
		}
	}
	for _, p := range ev.Removed {
		action := h.matchAction(p)
		if action == "ignore" {
			continue
		}
		e := alert.NewEvent("closed", p)
		if err := h.notifier.Send(e); err != nil {
			log.Printf("alert send error: %v", err)
		}
	}
}

func (h *WatcherHandler) matchAction(p ports.PortState) string {
	if h.matcher == nil {
		return "alert"
	}
	rule, ok := h.matcher.Match(p)
	if !ok {
		return "alert"
	}
	return rule.Action
}
