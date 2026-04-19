// Package monitor provides runtime monitoring components for portwatch.
//
// RetryNotifier
//
// RetryNotifier wraps any alert.Notifier and transparently retries delivery
// when the downstream notifier returns an error.  This is useful for
// transient failures such as network timeouts when using WebhookNotifier.
//
// Example usage:
//
//	webhook := monitor.NewWebhookNotifier(url, nil, nil)
//	retrying := monitor.NewRetryNotifier(webhook, 3, 500*time.Millisecond, logger)
//
// The notifier will attempt delivery up to maxTries times, sleeping delay
// between each attempt.  If all attempts fail the last error is returned to
// the caller.
package monitor
