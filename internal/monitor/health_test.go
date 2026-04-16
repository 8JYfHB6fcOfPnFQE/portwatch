package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHealthTracker_InitiallyHealthy(t *testing.T) {
	ht := NewHealthTracker()
	s := ht.Status()
	if !s.Healthy {
		t.Fatal("expected healthy on init")
	}
	if s.ConsecFails != 0 {
		t.Fatalf("expected 0 failures, got %d", s.ConsecFails)
	}
}

func TestHealthTracker_RecordFailure(t *testing.T) {
	ht := NewHealthTracker()
	ht.RecordFailure(errors.New("scan error"))
	s := ht.Status()
	if s.Healthy {
		t.Fatal("expected unhealthy after failure")
	}
	if s.ConsecFails != 1 {
		t.Fatalf("expected 1, got %d", s.ConsecFails)
	}
	if s.LastError != "scan error" {
		t.Fatalf("unexpected error: %s", s.LastError)
	}
}

func TestHealthTracker_RecordSuccess_Resets(t *testing.T) {
	ht := NewHealthTracker()
	ht.RecordFailure(errors.New("boom"))
	ht.RecordSuccess()
	s := ht.Status()
	if !s.Healthy {
		t.Fatal("expected healthy after success")
	}
	if s.ConsecFails != 0 {
		t.Fatalf("expected 0, got %d", s.ConsecFails)
	}
}

func TestHealthServer_HealthyResponse(t *testing.T) {
	ht := NewHealthTracker()
	ht.RecordSuccess()
	addr := "127.0.0.1:0"
	// Use a free port via net/http/httptest approach via direct server.
	// We pick a fixed high port for simplicity in tests.
	port := 19091
	addr = fmt.Sprintf("127.0.0.1:%d", port)
	hs := NewHealthServer(addr, ht)
	if err := hs.Start(); err != nil {
		t.Skipf("port unavailable, skipping: %v", err)
	}
	defer hs.Stop(context.Background()) //nolint:errcheck

	time.Sleep(20 * time.Millisecond)
	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var status HealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !status.Healthy {
		t.Fatal("expected healthy in response")
	}
}
