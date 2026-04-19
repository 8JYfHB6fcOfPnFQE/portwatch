package monitor_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func makeTeamsEvent() alert.Event {
	return alert.NewEvent(8080, "tcp", "opened", time.Now())
}

func TestTeamsNotifier_PostsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := monitor.NewTeamsNotifier(ts.URL, nil)
	if err := n.Send(makeTeamsEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["@type"] != "MessageCard" {
		t.Errorf("expected @type MessageCard, got %v", received["@type"])
	}
	if received["summary"] == "" {
		t.Error("expected non-empty summary")
	}
}

func TestTeamsNotifier_ForwardsToNext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	forwarded := false
	next := alert.NotifierFunc(func(ev alert.Event) error {
		forwarded = true
		return nil
	})

	n := monitor.NewTeamsNotifier(ts.URL, next)
	if err := n.Send(makeTeamsEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !forwarded {
		t.Error("expected event to be forwarded to next")
	}
}

func TestTeamsNotifier_BadStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := monitor.NewTeamsNotifier(ts.URL, nil)
	if err := n.Send(makeTeamsEvent()); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestTeamsNotifier_BadURL_DoesNotPanic(t *testing.T) {
	n := monitor.NewTeamsNotifier("http://127.0.0.1:0/no-server", nil)
	if err := n.Send(makeTeamsEvent()); err == nil {
		t.Error("expected error for unreachable URL")
	}
}
