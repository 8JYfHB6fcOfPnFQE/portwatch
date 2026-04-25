package monitor

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// OpsNotifier wraps a notifier and records per-event send latency and
// success/failure counts to a MetricsCollector.
type OpsNotifier struct {
	next    alert.Notifier
	metrics *OpsMetrics
	logger  *log.Logger
}

// OpsMetrics holds lightweight operational counters for a notifier chain.
type OpsMetrics struct {
	Sent    int64
	Failed  int64
	TotalMs int64
}

// NewOpsNotifier wraps next with operational instrumentation.
// logger may be nil; a default logger is used in that case.
func NewOpsNotifier(next alert.Notifier, metrics *OpsMetrics, logger *log.Logger) *OpsNotifier {
	if logger == nil {
		logger = log.Default()
	}
	if metrics == nil {
		metrics = &OpsMetrics{}
	}
	return &OpsNotifier{next: next, metrics: metrics, logger: logger}
}

// Send records latency and delegates to the wrapped notifier.
func (o *OpsNotifier) Send(ctx context.Context, ev alert.Event) error {
	start := time.Now()
	err := o.next.Send(ctx, ev)
	elapsed := time.Since(start).Milliseconds()
	o.metrics.TotalMs += elapsed
	if err != nil {
		o.metrics.Failed++
		o.logger.Printf("[ops] send failed port=%d proto=%s err=%v latency_ms=%d",
			ev.Port, ev.Proto, err, elapsed)
		return err
	}
	o.metrics.Sent++
	return nil
}

// AvgLatencyMs returns the mean send latency in milliseconds.
// Returns 0 when no events have been sent.
func (o *OpsNotifier) AvgLatencyMs() float64 {
	total := o.metrics.Sent + o.metrics.Failed
	if total == 0 {
		return 0
	}
	return float64(o.metrics.TotalMs) / float64(total)
}

// SuccessRate returns the fraction of send attempts that succeeded as a value
// in [0, 1]. Returns 0 when no events have been attempted.
func (o *OpsNotifier) SuccessRate() float64 {
	total := o.metrics.Sent + o.metrics.Failed
	if total == 0 {
		return 0
	}
	return float64(o.metrics.Sent) / float64(total)
}
