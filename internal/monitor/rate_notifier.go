package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// RateNotifier limits the number of alerts forwarded per time window.
// Events exceeding the limit are dropped and a summary is logged.
type RateNotifier struct {
	next    alert.Notifier
	limit   int
	window  time.Duration
	mu      sync.Mutex
	count   int
	windowStart time.Time
	dropped int
}

// NewRateNotifier returns a RateNotifier that forwards at most limit events
// per window duration to next. A limit <= 0 disables rate limiting.
func NewRateNotifier(limit int, window time.Duration, next alert.Notifier) *RateNotifier {
	return &RateNotifier{
		next:        next,
		limit:       limit,
		window:      window,
		windowStart: time.Now(),
	}
}

// Send forwards the event if within the rate limit, otherwise drops it.
func (r *RateNotifier) Send(e alert.Event) error {
	if r.limit <= 0 {
		if r.next != nil {
			return r.next.Send(e)
		}
		return nil
	}

	r.mu.Lock()
	now := time.Now()
	if now.Sub(r.windowStart) >= r.window {
		if r.dropped > 0 {
			_ = r.emitDropSummary(r.dropped)
		}
		r.count = 0
		dropped := r.dropped
		r.dropped = 0
		r.windowStart = now
		_ = dropped
	}

	if r.count >= r.limit {
		r.dropped++
		r.mu.Unlock()
		return nil
	}
	r.count++
	r.mu.Unlock()

	if r.next != nil {
		return r.next.Send(e)
	}
	return nil
}

// Dropped returns the number of events dropped in the current window.
func (r *RateNotifier) Dropped() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.dropped
}

func (r *RateNotifier) emitDropSummary(n int) error {
	if r.next == nil || n == 0 {
		return nil
	}
	summary := alert.NewEvent("rate-limit", 0, "tcp")
	summary.Message = fmt.Sprintf("rate limiter dropped %d events in the last window", n)
	return r.next.Send(summary)
}
