package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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

func TestLoadJSON(t *testing.T) {
	path := writeTempFile(t, "config.json", `{"port": 8081, "host": "example.com", "debug": true}`)

	var cfg testAppConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 8081 || cfg.Host != "example.com" || cfg.Debug != true {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadJSONWithEnvOverride(t *testing.T) {
	path := writeTempFile(t, "config.json", `{"port": 8081, "host": "example.com"}`)

	t.Setenv("TEST_APP_PORT", "7070")

	var cfg testAppConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 7070 {
		t.Errorf("Port = %d, want 7070 (env override)", cfg.Port)
	}
	if cfg.Host != "example.com" {
		t.Errorf("Host = %q, want example.com (from file)", cfg.Host)
	}
}

func TestLoadFormatExplicit(t *testing.T) {
	path := writeTempFile(t, "config.cfg", `{"port": 1234, "host": "explicit"}`)

	var cfg testAppConfig
	if err := LoadFormat(path, FormatJSON, &cfg); err != nil {
		t.Fatalf("LoadFormat returned error: %v", err)
	}

	if cfg.Port != 1234 || cfg.Host != "explicit" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadUnrecognizedExtension(t *testing.T) {
	path := writeTempFile(t, "config.cfg", "port: 8080\n")

	var cfg testAppConfig
	if err := Load(path, &cfg); err == nil {
		t.Fatal("expected error for unrecognized extension, got nil")
	}
}

type testDBConfig struct {
	Host string `yaml:"host" json:"host" env:"TEST_DB_HOST"`
	Port int    `yaml:"port" json:"port" env:"TEST_DB_PORT"`
}

type testNestedConfig struct {
	Name string       `yaml:"name" json:"name" env:"TEST_APP_NAME"`
	DB   testDBConfig `yaml:"db" json:"db"`
}

func TestLoadNestedStructEnvOverride(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "name: myapp\ndb:\n  host: dbhost\n  port: 5432\n")

	t.Setenv("TEST_DB_PORT", "6543")

	var cfg testNestedConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Name != "myapp" {
		t.Errorf("Name = %q, want myapp", cfg.Name)
	}
	if cfg.DB.Host != "dbhost" {
		t.Errorf("DB.Host = %q, want dbhost (from file)", cfg.DB.Host)
	}
	if cfg.DB.Port != 6543 {
		t.Errorf("DB.Port = %d, want 6543 (env override)", cfg.DB.Port)
	}
}

type testSliceDurationConfig struct {
	Tags    []string      `yaml:"tags" env:"TEST_TAGS"`
	Timeout time.Duration `yaml:"timeout" env:"TEST_TIMEOUT"`
}

func TestLoadSliceAndDurationFromFile(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "tags: [alpha, beta]\ntimeout: 5s\n")

	var cfg testSliceDurationConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if !reflect.DeepEqual(cfg.Tags, []string{"alpha", "beta"}) {
		t.Errorf("Tags = %v, want [alpha beta]", cfg.Tags)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", cfg.Timeout)
	}
}

func TestLoadSliceAndDurationFromEnv(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "tags: [alpha]\ntimeout: 5s\n")

	t.Setenv("TEST_TAGS", "x,y,z")
	t.Setenv("TEST_TIMEOUT", "2m")

	var cfg testSliceDurationConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if !reflect.DeepEqual(cfg.Tags, []string{"x", "y", "z"}) {
		t.Errorf("Tags = %v, want [x y z]", cfg.Tags)
	}
	if cfg.Timeout != 2*time.Minute {
		t.Errorf("Timeout = %v, want 2m", cfg.Timeout)
	}
}

func TestLoadFormatUnsupportedFormat(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\n")

	var cfg testAppConfig
	if err := LoadFormat(path, Format("xml"), &cfg); err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestLoadEmptyFile(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "")

	var cfg testAppConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error for empty file: %v", err)
	}
	if cfg.Port != 0 || cfg.Host != "" {
		t.Fatalf("expected zero-value config for empty file, got %+v", cfg)
	}
}

func TestLoadEmptyJSONObject(t *testing.T) {
	path := writeTempFile(t, "config.json", "{}")

	var cfg testAppConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error for empty JSON object: %v", err)
	}
	if cfg.Port != 0 || cfg.Host != "" {
		t.Fatalf("expected zero-value config for empty JSON object, got %+v", cfg)
	}
}

func TestLoadInvalidYAMLSyntax(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: [unterminated\n")

	var cfg testAppConfig
	if err := Load(path, &cfg); err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}

func TestLoadInvalidJSONSyntax(t *testing.T) {
	path := writeTempFile(t, "config.json", "{not valid json")

	var cfg testAppConfig
	if err := Load(path, &cfg); err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestLoadEnvOverrideUnsupportedKindErrors(t *testing.T) {
	type cfgType struct {
		Data map[string]string `yaml:"data" env:"TEST_UNSUPPORTED_KIND"`
	}

	path := writeTempFile(t, "config.yaml", "data: {}\n")
	t.Setenv("TEST_UNSUPPORTED_KIND", "x")

	var cfg cfgType
	if err := Load(path, &cfg); err == nil {
		t.Fatal("expected error for env override of unsupported field kind, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	var cfg testAppConfig
	if err := Load(filepath.Join(t.TempDir(), "missing.yaml"), &cfg); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadMissingFileWrapsErrNotExist(t *testing.T) {
	var cfg testAppConfig
	err := Load(filepath.Join(t.TempDir(), "missing.yaml"), &cfg)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected error wrapping os.ErrNotExist, got: %v", err)
	}
}

func TestLoadRequiresPointer(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\n")

	var cfg testAppConfig
	if err := Load(path, cfg); err == nil {
		t.Fatal("expected error when out is not a pointer, got nil")
	}
}
