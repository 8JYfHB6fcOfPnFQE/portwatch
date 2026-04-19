package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// teamsPayload is the adaptive card message format for Microsoft Teams incoming webhooks.
type teamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	ThemeColor string         `json:"themeColor"`
	Summary    string         `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle    string      `json:"activityTitle"`
	ActivitySubtitle string      `json:"activitySubtitle"`
	Facts            []teamsFact `json:"facts"`
}

type teamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewTeamsNotifier posts an alert to a Microsoft Teams channel via an
// incoming webhook URL and then forwards the event to next (if non-nil).
func NewTeamsNotifier(webhookURL string, next alert.Notifier) alert.Notifier {
	return alert.NotifierFunc(func(ev alert.Event) error {
		color := "d9534f"
		if ev.Kind == "closed" {
			color = "5bc0de"
		}

		payload := teamsPayload{
			Type:       "MessageCard",
			Context:    "http://schema.org/extensions",
			ThemeColor: color,
			Summary:    fmt.Sprintf("portwatch: port %d/%s %s", ev.Port, ev.Proto, ev.Kind),
			Sections: []teamsSection{
				{
					ActivityTitle:    "portwatch alert",
					ActivitySubtitle: fmt.Sprintf("Port %d/%s was %s", ev.Port, ev.Proto, ev.Kind),
					Facts: []teamsFact{
						{Name: "Port", Value: fmt.Sprintf("%d", ev.Port)},
						{Name: "Protocol", Value: ev.Proto},
						{Name: "Event", Value: ev.Kind},
						{Name: "Time", Value: ev.Time.Format("2006-01-02 15:04:05 UTC")},
					},
				},
			},
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("teams: marshal: %w", err)
		}

		resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("teams: post: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
		}

		if next != nil {
			return next.Send(ev)
		}
		return nil
	})
}
