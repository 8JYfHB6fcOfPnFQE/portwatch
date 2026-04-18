package monitor

import (
	"sync"
	"time"
)

// SilenceRule suppresses alerts for a specific port/proto pair until a deadline.
type SilenceRule struct {
	Port     int
	Proto    string
	Deadline time.Time
}

// SilenceStore manages active silence rules.
type SilenceStore struct {
	mu      sync.RWMutex
	rules   []SilenceRule
	nowFunc func() time.Time
}

// NewSilenceStore returns an initialised SilenceStore.
func NewSilenceStore() *SilenceStore {
	return &SilenceStore{nowFunc: time.Now}
}

// Add registers a silence for the given port/proto for the given duration.
func (s *SilenceStore) Add(port int, proto string, dur time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = append(s.rules, SilenceRule{
		Port:     port,
		Proto:    proto,
		Deadline: s.nowFunc().Add(dur),
	})
}

// IsSilenced reports whether the given port/proto is currently silenced.
func (s *SilenceStore) IsSilenced(port int, proto string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := s.nowFunc()
	for _, r := range s.rules {
		if r.Port == port && r.Proto == proto && now.Before(r.Deadline) {
			return true
		}
	}
	return false
}

// Purge removes expired silence rules.
func (s *SilenceStore) Purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.nowFunc()
	kept := s.rules[:0]
	for _, r := range s.rules {
		if now.Before(r.Deadline) {
			kept = append(kept, r)
		}
	}
	s.rules = kept
}

// List returns a copy of all active silence rules.
func (s *SilenceStore) List() []SilenceRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SilenceRule, len(s.rules))
	copy(out, s.rules)
	return out
}
