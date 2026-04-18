package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// WebhookConfig holds configuration for the webhook notifier.
type WebhookConfig struct {
	URL     string
	Timeout time.Duration
	Headers map[string]string
}

// webhookNotifier sends alert events as JSON POST requests to a webhook URL.
type webhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
	next   alert.Notifier
}

// NewWebhookNotifier returns a Notifier that POSTs events to a webhook and
// forwards to next regardless of HTTP outcome.
func NewWebhookNotifier(cfg WebhookConfig, next alert.Notifier) alert.Notifier {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	return &webhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
		next:   next,
	}
}

func (w *webhookNotifier) Send(e alert.Event) error {
	payload, err := json.Marshal(map[string]interface{}{
		"port":      e.Port,
		"proto":     e.Proto,
		"change":    e.Change,
		"timestamp": e.Timestamp.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.cfg.URL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		// non-fatal: log and continue to next
		fmt.Printf("webhook: send error: %v\n", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode >= 400 {
			fmt.Printf("webhook: unexpected status %d\n", resp.StatusCode)
		}
	}

	if w.next != nil {
		return w.next.Send(e)
	}
	return nil
}
