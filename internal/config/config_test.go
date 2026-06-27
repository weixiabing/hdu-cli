package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.Endpoint != "http://192.168.112.30" {
		t.Fatalf("expected default endpoint, got %q", cfg.Endpoint)
	}
	if cfg.CheckIntervalSeconds != 60 {
		t.Fatalf("expected 60 second interval, got %d", cfg.CheckIntervalSeconds)
	}
	if !cfg.AutoReconnect {
		t.Fatalf("expected auto reconnect to be enabled")
	}
}
