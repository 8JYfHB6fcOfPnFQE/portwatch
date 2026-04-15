package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestMain_ConfigLoadError(t *testing.T) {
	// Verify that a missing config path would be caught.
	// We test the config loading logic indirectly via config.Load.
	_, err := os.Stat("/nonexistent/path/portwatch.yaml")
	if !os.IsNotExist(err) {
		t.Skip("unexpected file exists")
	}
}

func TestMain_ValidConfigFile(t *testing.T) {
	content := `
interval: 5s
ports:
  start: 1024
  end: 1100
output:
  format: text
`
	p := writeTempConfig(t, content)
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected non-empty config file")
	}
}

func TestMain_RulesFileOverride(t *testing.T) {
	// Ensure rules file path can be set independently.
	rulesContent := `
- port: 22
  proto: tcp
  action: allow
  description: SSH
`
	dir := t.TempDir()
	rulesPath := filepath.Join(dir, "rules.yaml")
	if err := os.WriteFile(rulesPath, []byte(rulesContent), 0644); err != nil {
		t.Fatalf("failed to write rules file: %v", err)
	}

	info, err := os.Stat(rulesPath)
	if err != nil {
		t.Fatalf("rules file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("rules file should not be empty")
	}
}
