package monitor

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeSyslogEvent() alert.Event {
	return alert.NewEvent("opened", "tcp", 9090, "127.0.0.1")
}

func TestFormatSyslogLine_ContainsFields(t *testing.T) {
	e := makeSyslogEvent()
	line := formatSyslogLine(e)

	for _, want := range []string{"portwatch:", "action=opened", "proto=tcp", "port=9090", "addr=127.0.0.1"} {
		if !strings.Contains(line, want) {
			t.Errorf("expected %q in %q", want, line)
		}
	}
}

func TestFormatSyslogLine_NoAddr(t *testing.T) {
	e := alert.NewEvent("closed", "udp", 53, "")
	line := formatSyslogLine(e)
	if strings.Contains(line, "addr=") {
		t.Errorf("did not expect addr= in %q", line)
	}
	if !strings.Contains(line, "port=53") {
		t.Errorf("expected port=53 in %q", line)
	}
}

func TestSyslogNotifier_ForwardsToNext(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(e alert.Event) error {
		got = e
		return nil
	})

	// Use a real syslog writer; skip if unavailable (e.g., CI without syslogd).
	sn, err := NewSyslogNotifier(0, "portwatch-test", next)
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer sn.Close()

	e := makeSyslogEvent()
	if err := sn.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.Port != e.Port {
		t.Errorf("forwarded port = %d, want %d", got.Port, e.Port)
	}
}

func TestSyslogNotifier_NilNext_NoError(t *testing.T) {
	sn, err := NewSyslogNotifier(0, "portwatch-test", nil)
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer sn.Close()

	if err := sn.Send(makeSyslogEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
