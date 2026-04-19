package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// RetryNotifier wraps a Notifier and retries delivery on failure.
type RetryNotifier struct {
	next     alert.Notifier
	maxTries int
	delay    time.Duration
	log      *log.Logger
}

// NewRetryNotifier creates a RetryNotifier that will attempt delivery up to
// maxTries times, waiting delay between each attempt.
func NewRetryNotifier(next alert.Notifier, maxTries int, delay time.Duration, logger *log.Logger) *RetryNotifier {
	if maxTries < 1 {
		maxTries = 1
	}
	return &RetryNotifier{
		next:     next,
		maxTries: maxTries,
		delay:    delay,
		log:      logger,
	}
}

// Send attempts to deliver the event, retrying on error.
func (r *RetryNotifier) Send(e alert.Event) error {
	var err error
	for attempt := 1; attempt <= r.maxTries; attempt++ {
		if err = r.next.Send(e); err == nil {
			return nil
		}
		if r.log != nil {
			r.log.Printf("retry_notifier: attempt %d/%d failed: %v", attempt, r.maxTries, err)
		}
		if attempt < r.maxTries {
			time.Sleep(r.delay)
		}
	}
	return err
}
