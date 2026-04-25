package config

import "testing"

func validRoutingRule() RoutingRuleConfig {
	return RoutingRuleConfig{Field: "proto", Value: "tcp", Destination: "slack"}
}

func TestRoutingConfig_Validate_Disabled(t *testing.T) {
	rc := RoutingConfig{Enabled: false, Rules: []RoutingRuleConfig{
		{Field: "", Value: "", Destination: ""},
	}}
	if err := rc.Validate(); err != nil {
		t.Fatalf("disabled config should not error: %v", err)
	}
}

func TestRoutingConfig_Validate_Valid(t *testing.T) {
	rc := RoutingConfig{
		Enabled: true,
		Rules:   []RoutingRuleConfig{validRoutingRule()},
	}
	if err := rc.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRoutingConfig_Validate_MissingField(t *testing.T) {
	rc := RoutingConfig{
		Enabled: true,
		Rules:   []RoutingRuleConfig{{Field: "", Value: "tcp", Destination: "slack"}},
	}
	if err := rc.Validate(); err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestRoutingConfig_Validate_MissingValue(t *testing.T) {
	rc := RoutingConfig{
		Enabled: true,
		Rules:   []RoutingRuleConfig{{Field: "proto", Value: "", Destination: "slack"}},
	}
	if err := rc.Validate(); err == nil {
		t.Fatal("expected error for missing value")
	}
}

func TestRoutingConfig_Validate_MissingDestination(t *testing.T) {
	rc := RoutingConfig{
		Enabled: true,
		Rules:   []RoutingRuleConfig{{Field: "proto", Value: "tcp", Destination: ""}},
	}
	if err := rc.Validate(); err == nil {
		t.Fatal("expected error for missing destination")
	}
}

func TestRoutingConfig_Default_Disabled(t *testing.T) {
	rc := DefaultRoutingConfig()
	if rc.Enabled {
		t.Fatal("default routing config should be disabled")
	}
	if len(rc.Rules) != 0 {
		t.Fatalf("expected no default rules, got %d", len(rc.Rules))
	}
}
