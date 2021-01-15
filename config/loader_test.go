package config

import (
	"os"
	"path/filepath"
	"testing"
)

type testAppConfig struct {
	Port  int    `yaml:"port" env:"TEST_APP_PORT"`
	Host  string `yaml:"host" env:"TEST_APP_HOST"`
	Debug bool   `yaml:"debug" env:"TEST_APP_DEBUG"`
}

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestLoadYAML(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: localhost\ndebug: false\n")

	var cfg testAppConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 8080 || cfg.Host != "localhost" || cfg.Debug != false {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadYAMLWithEnvOverride(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: localhost\ndebug: false\n")

	t.Setenv("TEST_APP_PORT", "9090")
	t.Setenv("TEST_APP_DEBUG", "true")

	var cfg testAppConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090 (env override)", cfg.Port)
	}
	if cfg.Host != "localhost" {
		t.Errorf("Host = %q, want localhost (from file)", cfg.Host)
	}
	if cfg.Debug != true {
		t.Errorf("Debug = %v, want true (env override)", cfg.Debug)
	}
}

func TestLoadMissingFile(t *testing.T) {
	var cfg testAppConfig
	if err := Load(filepath.Join(t.TempDir(), "missing.yaml"), &cfg); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadRequiresPointer(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\n")

	var cfg testAppConfig
	if err := Load(path, cfg); err == nil {
		t.Fatal("expected error when out is not a pointer, got nil")
	}
}
