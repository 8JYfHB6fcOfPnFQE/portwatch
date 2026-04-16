// Package ports — watch.go
//
// Watcher provides a high-level continuous monitoring loop built on top of
// Scanner and History. It emits ChangeEvent values on a channel whenever
// ports are opened or closed between successive scans.
//
// Usage:
//
//	 scanner, _ := ports.NewScanner(ports.DefaultPortRange)
//	 history  := ports.NewHistory()
//	 filter   := ports.NewFilter(cfg)
//	 watcher  := ports.NewWatcher(scanner, history, ports.WatchConfig{
//	     Interval: 5 * time.Second,
//	     Filter:   filter,
//	 })
//	 ctx, cancel := context.WithCancel(context.Background())
//	 defer cancel()
//	 for ev := range watcher.Watch(ctx) {
//	     fmt.Println("added:", ev.Added, "removed:", ev.Removed)
//	 }
package ports
