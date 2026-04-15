package alert

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Format controls how events are serialised.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// FormattedNotifier wraps Notifier with pluggable output formats.
type FormattedNotifier struct {
	notifier *Notifier
	format   Format
}

// NewFormattedNotifier creates a FormattedNotifier.
func NewFormattedNotifier(format Format, writers ...io.Writer) *FormattedNotifier {
	return &FormattedNotifier{
		notifier: NewNotifier(writers...),
		format:   format,
	}
}

// Send dispatches the event using the configured format.
func (f *FormattedNotifier) Send(e Event) {
	switch f.format {
	case FormatJSON:
		f.sendJSON(e)
	default:
		f.notifier.Send(e)
	}
}

type jsonEvent struct {
	Timestamp string `json:"timestamp"`
	Level     Level  `json:"level"`
	Port      int    `json:"port"`
	Proto     string `json:"proto"`
	Message   string `json:"message"`
}

func (f *FormattedNotifier) sendJSON(e Event) {
	je := jsonEvent{
		Timestamp: e.Timestamp.Format(time.RFC3339),
		Level:     e.Level,
		Port:      e.Port,
		Proto:     e.Proto,
		Message:   e.Message,
	}
	b, err := json.Marshal(je)
	if err != nil {
		fmt.Fprintf(f.notifier.writers[0], "alert: json marshal error: %v\n", err)
		return
	}
	for _, w := range f.notifier.writers {
		_, _ = fmt.Fprintf(w, "%s\n", b)
	}
}
