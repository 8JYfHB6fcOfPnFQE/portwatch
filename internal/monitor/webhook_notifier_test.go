package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeWebhookEvent() alert.Event {
	return alert.NewEvent(8080, "tcp", "opened")
}

func TestWebhookNotifier_PostsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(WebhookConfig{URL: ts.URL}, nil)
	if err := n.Send(makeWebhookEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["proto"] != "tcp" {
		t.Errorf("expected proto=tcp, got %v", received["proto"])
	}
	if received["change"] != "opened" {
		t.Errorf("expected change=opened, got %v", received["change"])
	}
}

func TestWebhookNotifier_ForwardsToNext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	forwarded := false
	next := alert.NotifierFunc(func(e alert.Event) error {
		forwarded = true
		return nil
	})

	n := NewWebhookNotifier(WebhookConfig{URL: ts.URL}, next)
	n.Send(makeWebhookEvent())

	if !forwarded {
		t.Error("expected event to be forwarded to next notifier")
	}
}

func TestWebhookNotifier_CustomHeaders(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := WebhookConfig{
		URL:     ts.URL,
		Headers: map[string]string{"Authorization": "Bearer secret"},
		Timeout: 2 * time.Second,
	}
	n := NewWebhookNotifier(cfg, nil)
	n.Send(makeWebhookEvent())

	if authHeader != "Bearer secret" {
		t.Errorf("expected auth header, got %q", authHeader)
	}
}

func TestWebhookNotifier_BadURL_DoesNotPanic(t *testing.T) {
	n := NewWebhookNotifier(WebhookConfig{URL: "http://127.0.0.1:0/no-server"}, nil)
	// should not panic or return a fatal error
	n.Send(makeWebhookEvent())
}
