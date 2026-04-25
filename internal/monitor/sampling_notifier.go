package monitor

import (
	"math/rand"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// SamplingNotifier forwards only a statistical sample of events to the next
// notifier. A rate of 1.0 forwards all events; 0.0 drops all events.
type SamplingNotifier struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
	next alert.Notifier
}

// NewSamplingNotifier creates a SamplingNotifier with the given sample rate
// (0.0–1.0). Events are forwarded with probability equal to rate. A rate
// outside [0, 1] is clamped to the nearest boundary.
func NewSamplingNotifier(rate float64, next alert.Notifier) *SamplingNotifier {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &SamplingNotifier{
		rate: rate,
		rng:  rand.New(rand.NewSource(rand.Int63())), //nolint:gosec
		next: next,
	}
}

// Send forwards the event to the next notifier with probability s.rate.
func (s *SamplingNotifier) Send(ev alert.Event) error {
	s.mu.Lock()
	sampled := s.rng.Float64() < s.rate
	s.mu.Unlock()

	if !sampled {
		return nil
	}
	if s.next == nil {
		return nil
	}
	return s.next.Send(ev)
}

// Rate returns the configured sample rate.
func (s *SamplingNotifier) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}
