package config

import "fmt"

// EmailConfig holds optional SMTP alert settings.
type EmailConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

// Validate checks that required fields are set when email is enabled.
func (e EmailConfig) Validate() error {
	if !e.Enabled {
		return nil
	}
	if e.Host == "" {
		return fmt.Errorf("email: host is required")
	}
	if e.Port <= 0 || e.Port > 65535 {
		return fmt.Errorf("email: invalid port %d", e.Port)
	}
	if e.From == "" {
		return fmt.Errorf("email: from address is required")
	}
	if len(e.To) == 0 {
		return fmt.Errorf("email: at least one recipient required")
	}
	return nil
}
