package monitor

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// DigestNotifier batches events over a window and forwards a summary to next.
type DigestNotifier struct {
	mu       sync.Mutex
	events   []alert.Event
	window   time.Duration
	next     alert.Notifier
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewDigestNotifier creates a DigestNotifier that flushes every window duration.
func NewDigestNotifier(window time.Duration, next alert.Notifier) *DigestNotifier {
	d := &DigestNotifier{
		window: window,
		next:   next,
		stopCh: make(chan struct{}),
	}
	d.wg.Add(1)
	go d.run()
	return d
}

func (d *DigestNotifier) Send(e alert.Event) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = append(d.events, e)
	return nil
}

func (d *DigestNotifier) Stop() {
	close(d.stopCh)
	d.wg.Wait()
}

func (d *DigestNotifier) run() {
	defer d.wg.Done()
	ticker := time.NewTicker(d.window)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.flush()
		case <-d.stopCh:
			d.flush()
			return
		}
	}
}

func (d *DigestNotifier) flush() {
	d.mu.Lock()
	events := d.events
	d.events = nil
	d.mu.Unlock()
	if len(events) == 0 {
		return
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "digest: %d event(s)\n", len(events))
	for _, e := range events {
		fmt.Fprintf(&sb, "  [%s] port=%d proto=%s\n", e.Kind, e.Port, e.Proto)
	}
	summary := alert.NewEvent("digest", 0, "")
	summary.Message = strings.TrimRight(sb.String(), "\n")
	_ = d.next.Send(summary)
}
