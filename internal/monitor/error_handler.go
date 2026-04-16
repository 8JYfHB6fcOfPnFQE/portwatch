package monitor

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// ErrorHandler wraps a Throttle and logs scan errors with back-off information.
type ErrorHandler struct {
	throttle *ports.Throttle
	out      io.Writer
}

// NewErrorHandler creates an ErrorHandler using the given throttle config.
func NewErrorHandler(cfg ports.ThrottleConfig) *ErrorHandler {
	return &ErrorHandler{
		throttle: ports.NewThrottle(cfg),
		out:      os.Stderr,
	}
}

// WithWriter replaces the output writer (useful for tests).
func (e *ErrorHandler) WithWriter(w io.Writer) *ErrorHandler {
	e.out = w
	return e
}

// Handle records the error, logs it, and returns the recommended back-off delay.
func (e *ErrorHandler) Handle(err error) time.Duration {
	delay := e.throttle.Failure()
	fmt.Fprintf(e.out, "[portwatch] scan error (#%d): %v — retrying in %v\n",
		e.throttle.Consecutive(), err, delay)
	return delay
}

// OK signals a successful scan, resetting the back-off state.
func (e *ErrorHandler) OK() {
	e.throttle.Success()
}

// Consecutive returns the current consecutive error count.
func (e *ErrorHandler) Consecutive() int {
	return e.throttle.Consecutive()
}
