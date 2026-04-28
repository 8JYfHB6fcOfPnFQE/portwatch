package monitor_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func makeEnvelopeEvent() alert.Event {
	return alert.NewEvent("opened", "tcp", 9090, "127.0.0.1:9090")
}

func TestEnvelopeNotifier_StampsSeq(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error { got = ev; return nil })
	n := monitor.NewEnvelopeNotifier("host1", next)

	_ = n.Send(makeEnvelopeEvent())

	if got.Meta["envelope.seq"] != "1" {
		t.Fatalf("expected seq=1, got %q", got.Meta["envelope.seq"])
	}
}

func TestEnvelopeNotifier_IncrementsSeq(t *testing.T) {
	var last alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error { last = ev; return nil })
	n := monitor.NewEnvelopeNotifier("", next)

	for i := 1; i <= 3; i++ {
		_ = n.Send(makeEnvelopeEvent())
		seq, _ := strconv.Atoi(last.Meta["envelope.seq"])
		if seq != i {
			t.Fatalf("iteration %d: expected seq=%d, got %d", i, i, seq)
		}
	}
}

func TestEnvelopeNotifier_StampsSentAt(t *testing.T) {
	before := time.Now().UTC()
	var got alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error { got = ev; return nil })
	n := monitor.NewEnvelopeNotifier("srv", next)
	_ = n.Send(makeEnvelopeEvent())

	ts, err := time.Parse(time.RFC3339Nano, got.Meta["envelope.sent_at"])
	if err != nil {
		t.Fatalf("invalid sent_at: %v", err)
	}
	if ts.Before(before) {
		t.Fatalf("sent_at %v is before test start %v", ts, before)
	}
}

func TestEnvelopeNotifier_StampsSource(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error { got = ev; return nil })
	n := monitor.NewEnvelopeNotifier("edge-node-42", next)
	_ = n.Send(makeEnvelopeEvent())

	if got.Meta["envelope.source"] != "edge-node-42" {
		t.Fatalf("expected source=edge-node-42, got %q", got.Meta["envelope.source"])
	}
}

func TestEnvelopeNotifier_EmptySource_OmitsKey(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error { got = ev; return nil })
	n := monitor.NewEnvelopeNotifier("", next)
	_ = n.Send(makeEnvelopeEvent())

	if _, ok := got.Meta["envelope.source"]; ok {
		t.Fatal("expected envelope.source to be absent when source is empty")
	}
}

func TestEnvelopeNotifier_NilNext_NoError(t *testing.T) {
	n := monitor.NewEnvelopeNotifier("x", nil)
	if err := n.Send(makeEnvelopeEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
