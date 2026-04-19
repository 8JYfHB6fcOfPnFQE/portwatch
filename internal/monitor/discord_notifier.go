package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// discordPayload is the JSON body sent to a Discord webhook.
type discordPayload struct {
	Content string `json:"content"`
}

// DiscordNotifier posts an alert message to a Discord webhook URL.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
	next       alert.Notifier
}

// NewDiscordNotifier creates a DiscordNotifier that posts to webhookURL and
// optionally forwards to next.
func NewDiscordNotifier(webhookURL string, client *http.Client, next alert.Notifier) *DiscordNotifier {
	if client == nil {
		client = http.DefaultClient
	}
	return &DiscordNotifier{webhookURL: webhookURL, client: client, next: next}
}

// Send posts the event to Discord then forwards to the next notifier.
func (d *DiscordNotifier) Send(e alert.Event) error {
	text := fmt.Sprintf("**portwatch** `%s` port %d/%s — %s",
		e.Change, e.Port, e.Proto, e.Addr)

	body, err := json.Marshal(discordPayload{Content: text})
	if err != nil {
		return fmt.Errorf("discord: marshal: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}

	if d.next != nil {
		return d.next.Send(e)
	}
	return nil
}
