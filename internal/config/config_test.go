package config

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestLoadMissingFileReturnsDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error for missing config, got %v", err)
	}

	want := Default()
	if cfg != want {
		t.Fatalf("expected default config %+v, got %+v", want, cfg)
	}
}

func TestLoadMergesPartialConfigWithDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("endpoint: https://example.com\nlogLevel: debug\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Endpoint != "https://example.com" {
		t.Fatalf("expected endpoint override, got %q", cfg.Endpoint)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("expected log level override, got %q", cfg.LogLevel)
	}
	if cfg.CheckIntervalSeconds != 60 {
		t.Fatalf("expected default interval to remain 60, got %d", cfg.CheckIntervalSeconds)
	}
	if !cfg.AutoReconnect {
		t.Fatalf("expected default auto reconnect to remain enabled")
	}
}

func TestLoadMapsLegacyNetConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.yaml")
	legacy := "net:\n  endpoint: http://legacy.example\n  acid: \"3\"\n  auth:\n    username: legacy-user\n"
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Endpoint != "http://legacy.example" {
		t.Fatalf("expected legacy endpoint mapping, got %q", cfg.Endpoint)
	}
	if cfg.ACID != "3" {
		t.Fatalf("expected legacy acid mapping, got %q", cfg.ACID)
	}
	if cfg.Username != "legacy-user" {
		t.Fatalf("expected legacy username mapping, got %q", cfg.Username)
	}
}

func TestLoadPrefersTopLevelKeysOverLegacyNetConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "conflict.yaml")
	conflict := "" +
		"endpoint: https://new.example\n" +
		"acid: \"9\"\n" +
		"username: new-user\n" +
		"net:\n" +
		"  endpoint: http://legacy.example\n" +
		"  acid: \"3\"\n" +
		"  auth:\n" +
		"    username: legacy-user\n"
	if err := os.WriteFile(path, []byte(conflict), 0o644); err != nil {
		t.Fatalf("write conflicting config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Endpoint != "https://new.example" {
		t.Fatalf("expected top-level endpoint to win, got %q", cfg.Endpoint)
	}
	if cfg.ACID != "9" {
		t.Fatalf("expected top-level acid to win, got %q", cfg.ACID)
	}
	if cfg.Username != "new-user" {
		t.Fatalf("expected top-level username to win, got %q", cfg.Username)
	}
}

func TestLoadInvalidYAMLReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.yaml")
	if err := os.WriteFile(path, []byte("endpoint: [\n"), 0o644); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatalf("expected invalid yaml error")
	}
}
