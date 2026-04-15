// Package alert provides notification primitives for portwatch.
//
// It defines an Event type that carries information about a port-state
// change, a Notifier that fans the event out to one or more io.Writer
// destinations, and a FormattedNotifier that supports plain-text and
// JSON output formats.
//
// Typical usage:
//
//	notifier := alert.NewFormattedNotifier(alert.FormatJSON, os.Stdout, logFile)
//	notifier.Send(alert.NewEvent(alert.LevelAlert, 8080, "tcp", "unexpected port opened"))
package alert
