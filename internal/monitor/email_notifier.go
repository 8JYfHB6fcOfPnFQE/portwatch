package monitor

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// EmailConfig holds SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// emailNotifier sends alert events via SMTP and forwards to next.
type emailNotifier struct {
	cfg  EmailConfig
	next alert.Notifier
}

// NewEmailNotifier returns a Notifier that emails on each event.
func NewEmailNotifier(cfg EmailConfig, next alert.Notifier) alert.Notifier {
	return &emailNotifier{cfg: cfg, next: next}
}

func (e *emailNotifier) Send(ev alert.Event) error {
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)

	subject := fmt.Sprintf("[portwatch] %s on port %d/%s", ev.Kind, ev.Port, ev.Proto)
	body := fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\n\r\n%s\r\n",
		subject,
		e.cfg.From,
		strings.Join(e.cfg.To, ", "),
		ev.String(),
	)

	if err := smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, []byte(body)); err != nil {
		return fmt.Errorf("email notifier: %w", err)
	}

	if e.next != nil {
		return e.next.Send(ev)
	}
	return nil
}
