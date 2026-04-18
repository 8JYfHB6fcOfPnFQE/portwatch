package monitor

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestSilenceStore_IsSilenced_Active(t *testing.T) {
	s := NewSilenceStore()
	s.Add(8080, "tcp", time.Minute)
	if !s.IsSilenced(8080, "tcp") {
		t.Fatal("expected port to be silenced")
	}
}

func TestSilenceStore_IsSilenced_Expired(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.nowFunc = func() time.Time { return now }
	s.Add(8080, "tcp", time.Millisecond)
	s.nowFunc = func() time.Time { return now.Add(time.Second) }
	if s.IsSilenced(8080, "tcp") {
		t.Fatal("expected silence to have expired")
	}
}

func TestSilenceStore_Purge(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.nowFunc = func() time.Time { return now }
	s.Add(80, "tcp", time.Millisecond)
	s.Add(443, "tcp", time.Hour)
	s.nowFunc = func() time.Time { return now.Add(time.Second) }
	s.Purge()
	if len(s.List()) != 1 {
		t.Fatalf("expected 1 rule after purge, got %d", len(s.List()))
	}
}

func TestSilenceStore_DifferentProto_NotSilenced(t *testing.T) {
	s := NewSilenceStore()
	s.Add(53, "tcp", time.Minute)
	if s.IsSilenced(53, "udp") {
		t.Fatal("udp should not be silenced when only tcp is")
	}
}

func TestSilenceNotifier_DropsWhenSilenced(t *testing.T) {
	var buf bytes.Buffer
	inner := alert.NewNotifier(&buf)
	store := NewSilenceStore()
	store.Add(9000, "tcp", time.Minute)
	sn := NewSilenceNotifier(inner, store)
	ev := alert.NewEvent("opened", 9000, "tcp")
	if err := sn.Send(ev); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestSilenceNotifier_ForwardsWhenNotSilenced(t *testing.T) {
	var buf bytes.Buffer
	inner := alert.NewNotifier(&buf)
	store := NewSilenceStore()
	sn := NewSilenceNotifier(inner, store)
	ev := alert.NewEvent("opened", 9000, "tcp")
	if err := sn.Send(ev); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected output but got none")
	}
}

func TestSilenceHandler_List_Empty(t *testing.T) {
	var buf bytes.Buffer
	h := NewSilenceHandler(NewSilenceStore(), &buf)
	h.List()
	if buf.String() != "no active silences\n" {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestSilenceHandler_Add_And_List(t *testing.T) {
	var buf bytes.Buffer
	h := NewSilenceHandler(NewSilenceStore(), &buf)
	h.Add(443, "tcp", 10*time.Minute)
	buf.Reset()
	h.List()
	if buf.Len() == 0 {
		t.Fatal("expected list output")
	}
}
