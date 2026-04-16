package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved set of expected open ports.
type Baseline struct {
	CreatedAt time.Time            `json:"created_at"`
	Ports     map[string]PortState `json:"ports"`
}

// BaselineStore persists and loads port baselines.
type BaselineStore struct {
	path string
}

// NewBaselineStore creates a BaselineStore backed by the given file path.
func NewBaselineStore(path string) *BaselineStore {
	return &BaselineStore{path: path}
}

// Save writes the current port states as the new baseline.
func (b *BaselineStore) Save(states []PortState) error {
	m := make(map[string]PortState, len(states))
	for _, ps := range states {
		m[key(ps.Port, ps.Proto)] = ps
	}
	bl := Baseline{CreatedAt: time.Now(), Ports: m}
	data, err := json.MarshalIndent(bl, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline marshal: %w", err)
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Load reads the baseline from disk. Returns nil, nil if the file does not exist.
func (b *BaselineStore) Load() (*Baseline, error) {
	data, err := os.ReadFile(b.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline read: %w", err)
	}
	var bl Baseline
	if err := json.Unmarshal(data, &bl); err != nil {
		return nil, fmt.Errorf("baseline unmarshal: %w", err)
	}
	return &bl, nil
}

// Diff returns ports that are in current but not in the baseline (new)
// and ports that are in the baseline but not in current (removed).
func (b *Baseline) Diff(current []PortState) (added, removed []PortState) {
	curMap := make(map[string]PortState, len(current))
	for _, ps := range current {
		curMap[key(ps.Port, ps.Proto)] = ps
	}
	for k, ps := range b.Ports {
		if _, ok := curMap[k]; !ok {
			removed = append(removed, ps)
		}
	}
	for k, ps := range curMap {
		if _, ok := b.Ports[k]; !ok {
			added = append(added, ps)
		}
	}
	return
}

// Path returns the file path used by the store.
func (b *BaselineStore) Path() string { return b.path }
