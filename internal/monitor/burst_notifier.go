package monitor

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// BurstNotifier suppresses alerts when more than MaxEvents occur within Window.
// Once the burst threshold is exceeded, events are dropped until the window
// resets. This protects downstream notifiers from alert storms.
type BurstNotifier struct {
	mu        sync.Mutex
	next      alert.Notifier
	maxEvents int
	window    time.Duration
	counts    map[string][]time.Time
}

// NewBurstNotifier returns a BurstNotifier that forwards to next.
// maxEvents is the maximum number of events allowed per key within window.
// A maxEvents value of 0 disables burst limiting.
func NewBurstNotifier(next alert.Notifier, maxEvents int, window time.Duration) *BurstNotifier {
	return &BurstNotifier{
		next:      next,
		maxEvents: maxEvents,
		window:    window,
		counts:    make(map[string][]time.Time),
	}
}

// Send forwards the event if the burst limit has not been exceeded for the
// event's port+proto key within the configured window.
func (b *BurstNotifier) Send(ev alert.Event) error {
	if b.maxEvents == 0 {
		if b.next != nil {
			return b.next.Send(ev)
		}
		return nil
	}

	key := itoa(ev.Port) + "/" + ev.Proto
	now := time.Now()

	b.mu.Lock()
	times := b.prune(b.counts[key], now)
	if len(times) >= b.maxEvents {
		b.mu.Unlock()
		return nil
	}
	b.counts[key] = append(times, now)
	b.mu.Unlock()

	if b.next != nil {
		return b.next.Send(ev)
	}
	return nil
}

// prune removes timestamps older than the window from ts.
func (b *BurstNotifier) prune(ts []time.Time, now time.Time) []time.Time {
	cutoff := now.Add(-b.window)
	i := 0
	for i < len(ts) && ts[i].Before(cutoff) {
		i++
	}
	return ts[i:]
}

// Reset clears all burst counters.
func (b *BurstNotifier) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.counts = make(map[string][]time.Time)
}
