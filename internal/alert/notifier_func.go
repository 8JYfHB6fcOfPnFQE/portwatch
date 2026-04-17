package alert

// NotifierFunc is a function adapter that implements Notifier.
type NotifierFunc func(Event) error

// Send calls the underlying function.
func (f NotifierFunc) Send(ev Event) error {
	return f(ev)
}
