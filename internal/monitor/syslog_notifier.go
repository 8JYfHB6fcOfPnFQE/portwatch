package monitor

import (
	"fmt"
	"log/syslog"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// SyslogNotifier writes alert events to the system syslog and forwards to next.
type SyslogNotifier struct {
	writer *syslog.Writer
	next   alert.Notifier
}

// NewSyslogNotifier creates a SyslogNotifier that writes to syslog with the
// given priority tag. next may be nil.
func NewSyslogNotifier(priority syslog.Priority, tag string, next alert.Notifier) (*SyslogNotifier, error) {
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open: %w", err)
	}
	return &SyslogNotifier{writer: w, next: next}, nil
}

// Send writes the event to syslog and forwards to next.
func (s *SyslogNotifier) Send(e alert.Event) error {
	msg := formatSyslogLine(e)
	if err := s.writer.Info(msg); err != nil {
		return fmt.Errorf("syslog: write: %w", err)
	}
	if s.next != nil {
		return s.next.Send(e)
	}
	return nil
}

// Close releases the syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}

func formatSyslogLine(e alert.Event) string {
	parts := []string{
		fmt.Sprintf("action=%s", e.Action),
		fmt.Sprintf("proto=%s", e.Proto),
		fmt.Sprintf("port=%d", e.Port),
	}
	if e.Addr != "" {
		parts = append(parts, fmt.Sprintf("addr=%s", e.Addr))
	}
	return "portwatch: " + strings.Join(parts, " ")
}
