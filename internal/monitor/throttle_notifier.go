package monitor

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// ThrottleNotifier wraps a Notifier and suppresses repeated alerts for the
// same port/proto pair within a configurable window.
type ThrottleNotifier struct {
	mu      sync.Mutex
	next    alert.Notifier
	window  time.Duration
	seen    map[string]time.Time
	nowFunc func() time.Time
}

// NewThrottleNotifier returns a ThrottleNotifier that forwards at most one
// alert per port/proto pair within window.
func NewThrottleNotifier(next alert.Notifier, window time.Duration) *ThrottleNotifier {
	return &ThrottleNotifier{
		next:    next,
		window:  window,
		seen:    make(map[string]time.Time),
		nowFunc: time.Now,
	}
}

// Send forwards the event only if the cooldown window has elapsed since the
// last alert for the same key.
func (t *ThrottleNotifier) Send(e alert.Event) error {
	key := e.Proto + ":" + itoa(e.Port)
	now := t.nowFunc()

	t.mu.Lock()
	last, exists := t.seen[key]
	if exists && now.Sub(last) < t.window {
		t.mu.Unlock()
		return nil
	}
	t.seen[key] = now
	t.mu.Unlock()

	return t.next.Send(e)
}

// Reset clears throttle state for a specific port/proto key.
func (t *ThrottleNotifier) Reset(proto string, port int) {
	key := proto + ":" + itoa(port)
	t.mu.Lock()
	delete(t.seen, key)
	t.mu.Unlock()
}
