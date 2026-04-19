package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// FileNotifier writes alert events to a log file in JSON or text format.
type FileNotifier struct {
	mu   sync.Mutex
	path string
	fmt  string // "json" or "text"
	next alert.Notifier
}

// NewFileNotifier creates a FileNotifier that appends events to path.
// format must be "json" or "text". next may be nil.
func NewFileNotifier(path, format string, next alert.Notifier) (*FileNotifier, error) {
	if format != "json" && format != "text" {
		return nil, fmt.Errorf("file_notifier: unsupported format %q", format)
	}
	// Ensure the file is writable/creatable.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("file_notifier: open %s: %w", path, err)
	}
	f.Close()
	return &FileNotifier{path: path, fmt: format, next: next}, nil
}

func (fn *FileNotifier) Send(e alert.Event) error {
	line, err := fn.format(e)
	if err != nil {
		return err
	}

	fn.mu.Lock()
	f, err := os.OpenFile(fn.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fn.mu.Unlock()
		return fmt.Errorf("file_notifier: open: %w", err)
	}
	_, werr := fmt.Fprintln(f, line)
	f.Close()
	fn.mu.Unlock()

	if werr != nil {
		return fmt.Errorf("file_notifier: write: %w", werr)
	}
	if fn.next != nil {
		return fn.next.Send(e)
	}
	return nil
}

func (fn *FileNotifier) format(e alert.Event) (string, error) {
	if fn.fmt == "json" {
		type record struct {
			Time  string `json:"time"`
			Kind  string `json:"kind"`
			Proto string `json:"proto"`
			Port  int    `json:"port"`
			Addr  string `json:"addr"`
		}{
		}
		r := struct {
			Time  string `json:"time"`
			Kind  string `json:"kind"`
			Proto string `json:"proto"`
			Port  int    `json:"port"`
			Addr  string `json:"addr"`
		}{
			Time:  e.Time.Format(time.RFC3339),
			Kind:  e.Kind,
			Proto: e.Proto,
			Port:  e.Port,
			Addr:  e.Addr,
		}
		b, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return fmt.Sprintf("%s %-8s %s:%d", e.Time.Format(time.RFC3339), e.Kind, e.Proto, e.Port), nil
}
