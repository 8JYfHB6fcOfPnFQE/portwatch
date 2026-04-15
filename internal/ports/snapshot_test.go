package ports

import (
	"os"
	"path/filepath"
	"testing"
)

func tempSnapshotPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestSnapshotStore_SaveAndLoad(t *testing.T) {
	store := NewSnapshotStore(tempSnapshotPath(t))

	states := map[string]PortState{
		"tcp:8080": {Port: 8080, Proto: "tcp", Open: true},
		"udp:53":   {Port: 53, Proto: "udp", Open: true},
	}

	if err := store.Save(states); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(loaded) != len(states) {
		t.Errorf("expected %d ports, got %d", len(states), len(loaded))
	}

	for k, want := range states {
		got, ok := loaded[k]
		if !ok {
			t.Errorf("missing key %q in loaded snapshot", k)
			continue
		}
		if got.Port != want.Port || got.Proto != want.Proto || got.Open != want.Open {
			t.Errorf("key %q: got %+v, want %+v", k, got, want)
		}
	}
}

func TestSnapshotStore_Load_MissingFile(t *testing.T) {
	store := NewSnapshotStore(tempSnapshotPath(t))

	states, err := store.Load()
	if err != nil {
		t.Fatalf("Load() on missing file returned error: %v", err)
	}
	if len(states) != 0 {
		t.Errorf("expected empty map, got %d entries", len(states))
	}
}

func TestSnapshotStore_Load_CorruptFile(t *testing.T) {
	path := tempSnapshotPath(t)
	if err := os.WriteFile(path, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	store := NewSnapshotStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error on corrupt JSON, got nil")
	}
}

func TestSnapshotStore_Path(t *testing.T) {
	path := "/tmp/test.json"
	store := NewSnapshotStore(path)
	if store.Path() != path {
		t.Errorf("Path() = %q, want %q", store.Path(), path)
	}
}
