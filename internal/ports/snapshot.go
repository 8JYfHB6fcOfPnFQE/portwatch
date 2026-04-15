package ports

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents a persisted view of port states at a point in time.
type Snapshot struct {
	Timestamp time.Time            `json:"timestamp"`
	Ports     map[string]PortState `json:"ports"`
}

// SnapshotStore handles reading and writing port snapshots to disk.
type SnapshotStore struct {
	path string
}

// NewSnapshotStore creates a SnapshotStore backed by the given file path.
func NewSnapshotStore(path string) *SnapshotStore {
	return &SnapshotStore{path: path}
}

// Save writes the current port states to the snapshot file.
func (s *SnapshotStore) Save(states map[string]PortState) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     states,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the snapshot file and returns the stored port states.
// Returns an empty map and no error if the file does not exist.
func (s *SnapshotStore) Load() (map[string]PortState, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return make(map[string]PortState), nil
	}
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return snap.Ports, nil
}

// Path returns the file path used by this store.
func (s *SnapshotStore) Path() string {
	return s.path
}
