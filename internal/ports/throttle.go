package ports

import "time"

// ThrottleConfig controls how the scanner backs off under repeated errors.
type ThrottleConfig struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
	Factor    float64
}

// DefaultThrottleConfig returns sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		BaseDelay: 500 * time.Millisecond,
		MaxDelay:  30 * time.Second,
		Factor:    2.0,
	}
}

// Throttle tracks consecutive errors and computes a back-off delay.
type Throttle struct {
	cfg          ThrottleConfig
	consecutive  int
	currentDelay time.Duration
}

// NewThrottle creates a Throttle with the given config.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	return &Throttle{cfg: cfg}
}

// Failure records a scan error and returns the delay to wait before retrying.
func (t *Throttle) Failure() time.Duration {
	t.consecutive++
	if t.consecutive == 1 {
		t.currentDelay = t.cfg.BaseDelay
	} else {
		t.currentDelay = time.Duration(float64(t.currentDelay) * t.cfg.Factor)
	}
	if t.currentDelay > t.cfg.MaxDelay {
		t.currentDelay = t.cfg.MaxDelay
	}
	return t.currentDelay
}

// Success resets the back-off state.
func (t *Throttle) Success() {
	t.consecutive = 0
	t.currentDelay = 0
}

// Consecutive returns the number of consecutive failures.
func (t *Throttle) Consecutive() int { return t.consecutive }
