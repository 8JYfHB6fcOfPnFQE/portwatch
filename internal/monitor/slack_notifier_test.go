package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeSlackEvent() alert.Event {
	return alert.NewEvent("opened", 9200, "tcp", time.Now())
}

func TestSlackNotifier_PostsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL, ts.Client(), nil)
	if err := n.Send(makeSlackEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text payload")
	}
}

func TestSlackNotifier_ForwardsToNext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	forwarded := false
	next := alert.NotifierFunc(func(e alert.Event) error {
		forwarded = true
		return nil
	})

	n := NewSlackNotifier(ts.URL, ts.Client(), next)
	if err := n.Send(makeSlackEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !forwarded {
		t.Error("expected event forwarded to next notifier")
	}
}

func TestSlackNotifier_BadStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL, ts.Client(), nil)
	if err := n.Send(makeSlackEvent()); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestSlackNotifier_BadURL_DoesNotPanic(t *testing.T) {
	n := NewSlackNotifier("http://127.0.0.1:0/no-server", nil, nil)
	if err := n.Send(makeSlackEvent()); err == nil {
		t.Error("expected error for unreachable URL")
	}
}
