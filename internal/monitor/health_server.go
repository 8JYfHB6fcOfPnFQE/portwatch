package monitor

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

// HealthServer exposes a simple HTTP endpoint for health checks.
type HealthServer struct {
	tracker *HealthTracker
	server  *http.Server
}

// NewHealthServer creates a HealthServer bound to addr (e.g. "127.0.0.1:9090").
func NewHealthServer(addr string, tracker *HealthTracker) *HealthServer {
	mux := http.NewServeMux()
	hs := &HealthServer{tracker: tracker}
	mux.HandleFunc("/healthz", hs.handleHealth)
	hs.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return hs
}

// Start begins listening in a background goroutine. It returns an error if the
// listener cannot be bound.
func (hs *HealthServer) Start() error {
	ln, err := net.Listen("tcp", hs.server.Addr)
	if err != nil {
		return err
	}
	go hs.server.Serve(ln) //nolint:errcheck
	return nil
}

// Stop gracefully shuts down the HTTP server.
func (hs *HealthServer) Stop(ctx context.Context) error {
	return hs.server.Shutdown(ctx)
}

func (hs *HealthServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := hs.tracker.Status()
	w.Header().Set("Content-Type", "application/json")
	if !status.Healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(status) //nolint:errcheck
}
