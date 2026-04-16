package ports

import (
	"testing"
	"time"
)

func TestThrottle_FirstFailure_ReturnsBaseDelay(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	d := th.Failure()
	if d != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", d)
	}
}

func TestThrottle_SecondFailure_Doubles(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	th.Failure()
	d := th.Failure()
	if d != time.Second {
		t.Fatalf("expected 1s, got %v", d)
	}
}

func TestThrottle_CapsAtMaxDelay(t *testing.T) {
	cfg := ThrottleConfig{BaseDelay: 10 * time.Second, MaxDelay: 30 * time.Second, Factor: 2.0}
	th := NewThrottle(cfg)
	var d time.Duration
	for i := 0; i < 10; i++ {
		d = th.Failure()
	}
	if d != 30*time.Second {
		t.Fatalf("expected max 30s, got %v", d)
	}
}

func TestThrottle_SuccessResetsState(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	th.Failure()
	th.Failure()
	th.Success()
	if th.Consecutive() != 0 {
		t.Fatalf("expected 0 consecutive, got %d", th.Consecutive())
	}
	d := th.Failure()
	if d != 500*time.Millisecond {
		t.Fatalf("expected base delay after reset, got %v", d)
	}
}

func TestThrottle_Consecutive_Tracks(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	th.Failure()
	th.Failure()
	th.Failure()
	if th.Consecutive() != 3 {
		t.Fatalf("expected 3, got %d", th.Consecutive())
	}
}
