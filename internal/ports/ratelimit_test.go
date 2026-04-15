package ports

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow_FirstCallPermitted(t *testing.T) {
	rl := NewRateLimiter(5 * time.Second)
	if !rl.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRateLimiter_Allow_SecondCallBlocked(t *testing.T) {
	rl := NewRateLimiter(5 * time.Second)
	rl.Allow("tcp:8080")
	if rl.Allow("tcp:8080") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestRateLimiter_Allow_DifferentKeysIndependent(t *testing.T) {
	rl := NewRateLimiter(5 * time.Second)
	rl.Allow("tcp:8080")
	if !rl.Allow("udp:53") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestRateLimiter_Allow_AfterCooldownExpires(t *testing.T) {
	now := time.Now()
	rl := NewRateLimiter(2 * time.Second)
	rl.now = func() time.Time { return now }

	rl.Allow("tcp:9090")

	// advance clock beyond cooldown
	rl.now = func() time.Time { return now.Add(3 * time.Second) }
	if !rl.Allow("tcp:9090") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestRateLimiter_Reset_ClearsKey(t *testing.T) {
	rl := NewRateLimiter(10 * time.Second)
	rl.Allow("tcp:443")
	rl.Reset("tcp:443")
	if !rl.Allow("tcp:443") {
		t.Fatal("expected call after Reset to be allowed")
	}
}

func TestRateLimiter_Len_TracksKeys(t *testing.T) {
	rl := NewRateLimiter(10 * time.Second)
	if rl.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", rl.Len())
	}
	rl.Allow("tcp:80")
	rl.Allow("tcp:443")
	if rl.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", rl.Len())
	}
}

func TestRateLimiter_Reset_UnknownKeyNoOp(t *testing.T) {
	rl := NewRateLimiter(5 * time.Second)
	// should not panic
	rl.Reset("tcp:9999")
	if rl.Len() != 0 {
		t.Fatalf("expected 0 keys after resetting unknown key, got %d", rl.Len())
	}
}
