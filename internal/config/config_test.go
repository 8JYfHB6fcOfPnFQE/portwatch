package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval)
	}
	if cfg.PortRange.From != 1 || cfg.PortRange.To != 65535 {
		t.Errorf("unexpected default port range: %+v", cfg.PortRange)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("expected text format, got %q", cfg.Output.Format)
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	path := writeTemp(t, `
scan_interval: 10s
port_range:
  from: 1024
  to: 9000
output:
  format: json
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.ScanInterval)
	}
	if cfg.PortRange.From != 1024 || cfg.PortRange.To != 9000 {
		t.Errorf("unexpected port range: %+v", cfg.PortRange)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected json, got %q", cfg.Output.Format)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_InvalidInterval(t *testing.T) {
	cfg := config.Default()
	cfg.ScanInterval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for zero interval")
	}
}

func TestValidate_InvalidPortRange(t *testing.T) {
	cfg := config.Default()
	cfg.PortRange.From = 9000
	cfg.PortRange.To = 1000
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for inverted range")
	}
}

func TestValidate_BadFormat(t *testing.T) {
	cfg := config.Default()
	cfg.Output.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for unsupported format")
	}
}
