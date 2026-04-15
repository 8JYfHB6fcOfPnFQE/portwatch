package ports

import (
	"sync"
	"time"
)

// RateLimiter suppresses repeated alerts for the same port/proto pair
// within a configurable cooldown window.
type RateLimiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// NewRateLimiter returns a RateLimiter with the given cooldown duration.
// Events for a given key are allowed at most once per cooldown window.
func NewRateLimiter(cooldown time.Duration) *RateLimiter {
	return &RateLimiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event identified by key should be forwarded.
// It returns false if an event with the same key was already allowed
// within the cooldown window.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	if t, ok := r.last[key]; ok && now.Sub(t) < r.cooldown {
		return false
	}
	r.last[key] = now
	return true
}

// Reset clears the rate-limit state for a specific key.
// Useful when a port transitions to a new state and the cooldown
// should restart.
func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.last, key)
}

// Len returns the number of tracked keys (for observability / testing).
func (r *RateLimiter) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.last)
}
