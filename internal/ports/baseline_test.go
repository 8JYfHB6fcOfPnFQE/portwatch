package ports

import (
	"os"
	"path/filepath"
	"testing"
)

func tempBaselinePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestBaselineStore_SaveAndLoad(t *testing.T) {
	store := NewBaselineStore(tempBaselinePath(t))
	states := []PortState{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	}
	if err := store.Save(states); err != nil {
		t.Fatalf("Save: %v", err)
	}
	bl, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(bl.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(bl.Ports))
	}
}

func TestBaselineStore_Load_MissingFile(t *testing.T) {
	store := NewBaselineStore("/nonexistent/baseline.json")
	bl, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bl != nil {
		t.Error("expected nil baseline for missing file")
	}
}

func TestBaselineStore_Load_CorruptFile(t *testing.T) {
	path := tempBaselinePath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	store := NewBaselineStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestBaseline_Diff_Added(t *testing.T) {
	bl := &Baseline{
		Ports: map[string]PortState{
			"80/tcp": {Port: 80, Proto: "tcp"},
		},
	}
	current := []PortState{
		{Port: 80, Proto: "tcp"},
		{Port: 9090, Proto: "tcp"},
	}
	added, removed := bl.Diff(current)
	if len(added) != 1 || added[0].Port != 9090 {
		t.Errorf("expected port 9090 added, got %v", added)
	}
	if len(removed) != 0 {
		t.Errorf("expected no removed ports, got %v", removed)
	}
}

func TestBaseline_Diff_Removed(t *testing.T) {
	bl := &Baseline{
		Ports: map[string]PortState{
			"80/tcp":   {Port: 80, Proto: "tcp"},
			"8080/tcp": {Port: 8080, Proto: "tcp"},
		},
	}
	current := []PortState{{Port: 80, Proto: "tcp"}}
	_, removed := bl.Diff(current)
	if len(removed) != 1 || removed[0].Port != 8080 {
		t.Errorf("expected port 8080 removed, got %v", removed)
	}
}

func TestBaselineStore_Path(t *testing.T) {
	store := NewBaselineStore("/tmp/bl.json")
	if store.Path() != "/tmp/bl.json" {
		t.Errorf("unexpected path: %s", store.Path())
	}
}
