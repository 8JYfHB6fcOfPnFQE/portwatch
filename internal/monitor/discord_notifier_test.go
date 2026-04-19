package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeDiscordEvent() alert.Event {
	return alert.NewEvent("opened", 9200, "tcp", "127.0.0.1", time.Now())
}

func TestDiscordNotifier_PostsJSON(t *testing.T) {
	var received discordPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	n := NewDiscordNotifier(srv.URL, srv.Client(), nil)
	if err := n.Send(makeDiscordEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestDiscordNotifier_ForwardsToNext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	var forwarded bool
	next := alert.NotifierFunc(func(e alert.Event) error { forwarded = true; return nil })
	n := NewDiscordNotifier(srv.URL, srv.Client(), next)
	n.Send(makeDiscordEvent())
	if !forwarded {
		t.Error("expected event forwarded to next")
	}
}

func TestDiscordNotifier_BadStatus_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := NewDiscordNotifier(srv.URL, srv.Client(), nil)
	if err := n.Send(makeDiscordEvent()); err == nil {
		t.Error("expected error on bad status")
	}
}

func TestDiscordNotifier_BadURL_DoesNotPanic(t *testing.T) {
	n := NewDiscordNotifier("http://127.0.0.1:0/nope", nil, nil)
	_ = n.Send(makeDiscordEvent())
}
