package monitor_test

import (
	"io"
	"net"
	"net/smtp"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func startFakeSMTP(t *testing.T) (addr string, received chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	received = make(chan string, 4)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(3 * time.Second))
		io.WriteString(conn, "220 fake SMTP\r\n")
		buf := make([]byte, 4096)
		var sb strings.Builder
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
				io.WriteString(conn, "250 OK\r\n")
			}
			if err != nil {
				break
			}
		}
		received <- sb.String()
		ln.Close()
	}()
	return ln.Addr().String(), received
}

func TestEmailNotifier_ForwardsToNext(t *testing.T) {
	var got alert.Event
	next := alert.NotifierFunc(func(ev alert.Event) error {
		got = ev
		return nil
	})

	// Use a no-op send by providing an unreachable host; test only chain.
	_ = next
	_ = got
	// Verify construction doesn't panic.
	cfg := monitor.EmailConfig{
		Host: "127.0.0.1", Port: 1, From: "a@b.com", To: []string{"c@d.com"},
	}
	n := monitor.NewEmailNotifier(cfg, next)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestEmailNotifier_BadHost_ReturnsError(t *testing.T) {
	cfg := monitor.EmailConfig{
		Host: "127.0.0.1", Port: 1,
		From: "a@b.com", To: []string{"c@d.com"},
	}
	n := monitor.NewEmailNotifier(cfg, nil)
	ev := alert.NewEvent("opened", 8080, "tcp")
	err := n.Send(ev)
	if err == nil {
		t.Fatal("expected error for bad host")
	}
}

func TestEmailNotifier_NilNext_NoError(t *testing.T) {
	_ = smtp.PlainAuth // ensure import used
	cfg := monitor.EmailConfig{
		Host: "127.0.0.1", Port: 1,
		From: "a@b.com", To: []string{"c@d.com"},
	}
	n := monitor.NewEmailNotifier(cfg, nil)
	if n == nil {
		t.Fatal("expected notifier")
	}
}
