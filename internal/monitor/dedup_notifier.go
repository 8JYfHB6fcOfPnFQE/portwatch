package monitor

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// DedupNotifier suppresses duplicate alerts for the same port/proto within a TTL window.
type DedupNotifier struct {
	inner alert.Notifier
	ttl   time.Duration
	mu    sync.Mutex
	seen  map[string]time.Time
}

// NewDedupNotifier wraps inner and drops events already seen within ttl.
func NewDedupNotifier(inner alert.Notifier, ttl time.Duration) *DedupNotifier {
	return &DedupNotifier{
		inner: inner,
		ttl:   ttl,
		seen:  make(map[string]time.Time),
	}
}

func (d *DedupNotifier) Send(e alert.Event) error {
	key := e.Proto + ":" + itoa(int(e.Port)) + ":" + e.Kind

	d.mu.Lock()
	if t, ok := d.seen[key]; ok && time.Since(t) < d.ttl {
		d.mu.Unlock()
		return nil
	}
	d.seen[key] = time.Now()
	d.mu.Unlock()

	return d.inner.Send(e)
}

// Flush removes all cached keys, useful for testing or forced resets.
func (d *DedupNotifier) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
