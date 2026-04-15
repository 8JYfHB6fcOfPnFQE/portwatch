package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change that triggered an alert.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Proto     string
	Message   string
}

// Notifier sends alert events to one or more outputs.
type Notifier struct {
	writers []io.Writer
}

// NewNotifier creates a Notifier that writes to the given writers.
// If no writers are provided, os.Stdout is used.
func NewNotifier(writers ...io.Writer) *Notifier {
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}
	return &Notifier{writers: writers}
}

// Send formats and dispatches the event to all configured writers.
func (n *Notifier) Send(e Event) {
	line := fmt.Sprintf("%s [%s] port=%d proto=%s msg=%q\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port,
		e.Proto,
		e.Message,
	)
	for _, w := range n.writers {
		_, _ = fmt.Fprint(w, line)
	}
}

// NewEvent is a convenience constructor for an Event.
func NewEvent(level Level, port int, proto, message string) Event {
	return Event{
		Timestamp: time.Now(),
		Level:     level,
		Port:      port,
		Proto:     proto,
		Message:   message,
	}
}
