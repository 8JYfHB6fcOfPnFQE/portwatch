package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// pagerDutyPayload is the V2 Events API request body.
type pagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     pagerDutyDetail   `json:"payload"`
	Client      string            `json:"client,omitempty"`
}

type pagerDutyDetail struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// PagerDutyNotifier sends alerts to PagerDuty and forwards to next.
type PagerDutyNotifier struct {
	routingKey string
	severity   string
	client     *http.Client
	next       alert.Notifier
}

// NewPagerDutyNotifier creates a PagerDutyNotifier.
// severity should be one of: critical, error, warning, info.
func NewPagerDutyNotifier(routingKey, severity string, next alert.Notifier) *PagerDutyNotifier {
	if severity == "" {
		severity = "error"
	}
	return &PagerDutyNotifier{
		routingKey: routingKey,
		severity:   severity,
		client:     &http.Client{Timeout: 10 * time.Second},
		next:       next,
	}
}

// Send posts the event to PagerDuty then forwards downstream.
func (p *PagerDutyNotifier) Send(e alert.Event) error {
	body := pagerDutyPayload{
		RoutingKey:  p.routingKey,
		EventAction: "trigger",
		Client:      "portwatch",
		Payload: pagerDutyDetail{
			Summary:   fmt.Sprintf("[portwatch] %s port %d/%s", e.Change, e.Port, e.Proto),
			Source:    "portwatch",
			Severity:  p.severity,
			Timestamp: e.Timestamp.UTC().Format(time.RFC3339),
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal: %w", err)
	}

	resp, err := p.client.Post(pagerDutyEventURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}

	if p.next != nil {
		return p.next.Send(e)
	}
	return nil
}
