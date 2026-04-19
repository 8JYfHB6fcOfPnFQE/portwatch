package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// SlackNotifier sends alert events to a Slack incoming webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
	next       alert.Notifier
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier returns a SlackNotifier that posts to webhookURL and
// forwards the event to next (may be nil).
func NewSlackNotifier(webhookURL string, client *http.Client, next alert.Notifier) *SlackNotifier {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &SlackNotifier{webhookURL: webhookURL, client: client, next: next}
}

// Send posts the event to Slack and forwards to the next notifier.
func (s *SlackNotifier) Send(e alert.Event) error {
	payload := slackPayload{Text: fmt.Sprintf("[portwatch] %s port %d/%s — %s",
		e.Kind, e.Port, e.Proto, e.Timestamp.Format(time.RFC3339))}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}

	if s.next != nil {
		return s.next.Send(e)
	}
	return nil
}
