package rules

import (
	"os"
	"testing"
)

func TestRule_Validate_Valid(t *testing.T) {
	r := Rule{Name: "ssh", Port: 22, Proto: "tcp", Action: ActionAlert}
	if err := r.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRule_Validate_InvalidPort(t *testing.T) {
	r := Rule{Name: "bad", Port: 0, Proto: "tcp", Action: ActionAlert}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for port 0, got nil")
	}
}

func TestRule_Validate_InvalidProto(t *testing.T) {
	r := Rule{Name: "bad", Port: 80, Proto: "icmp", Action: ActionAlert}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for proto 'icmp', got nil")
	}
}

func TestRule_Validate_InvalidAction(t *testing.T) {
	r := Rule{Name: "bad", Port: 80, Proto: "tcp", Action: "block"}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for action 'block', got nil")
	}
}

func TestMatcher_Match_Found(t *testing.T) {
	m, err := NewMatcher([]Rule{
		{Name: "http", Port: 80, Proto: "tcp", Action: ActionAlert},
		{Name: "dns", Port: 53, Proto: "udp", Action: ActionIgnore},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := m.Match(80, "tcp")
	if r == nil {
		t.Fatal("expected match for port 80/tcp, got nil")
	}
	if r.Name != "http" {
		t.Errorf("expected rule name 'http', got %q", r.Name)
	}
}

func TestMatcher_Match_NotFound(t *testing.T) {
	m, err := NewMatcher([]Rule{
		{Name: "ssh", Port: 22, Proto: "tcp", Action: ActionAlert},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r := m.Match(9999, "tcp"); r != nil {
		t.Errorf("expected no match, got rule %q", r.Name)
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	content := `rules:
  - name: ssh
    port: 22
    proto: tcp
    action: alert
  - name: dns
    port: 53
    proto: udp
    action: ignore
`
	f, err := os.CreateTemp("", "portwatch-rules-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	m, err := LoadFromFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r := m.Match(22, "tcp"); r == nil {
		t.Error("expected match for port 22/tcp")
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/rules.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
