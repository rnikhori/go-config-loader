package config

import (
	"reflect"
	"testing"
	"time"
)

type integrationDBConfig struct {
	Host string `yaml:"host" json:"host" env:"INTEG_DB_HOST" flag:"db-host" default:"db-default"`
	Port int    `yaml:"port" json:"port" env:"INTEG_DB_PORT" flag:"db-port" default:"5432" validate:"required"`
}

type integrationConfig struct {
	Name    string              `yaml:"name" json:"name" env:"INTEG_NAME" flag:"name" validate:"required"`
	Port    int                 `yaml:"port" json:"port" env:"INTEG_PORT" flag:"port" default:"8080"`
	Tags    []string            `yaml:"tags" json:"tags" env:"INTEG_TAGS" flag:"tags"`
	Timeout time.Duration       `yaml:"timeout" json:"timeout" env:"INTEG_TIMEOUT" flag:"timeout" default:"10s"`
	DB      integrationDBConfig `yaml:"db" json:"db"`
}

// TestFullPipelinePrecedence exercises every layer at once: a flag beats
// an env var, which beats a file value, which beats a default, and a
// field with none of those falls back to its default. It also checks
// that a nested field's flag override applies and that validation
// passes once every required field ends up set by some layer.
func TestFullPipelinePrecedence(t *testing.T) {
	path := writeTempFile(t, "config.yaml", `
name: file-name
port: 1111
tags: [file-tag]
db:
  host: file-db-host
  port: 2222
`)

	t.Setenv("INTEG_PORT", "2222")
	t.Setenv("INTEG_DB_PORT", "3333")
	t.Setenv("INTEG_TAGS", "env-a,env-b")

	args := []string{
		"-port=4444",
		"-db-host=flag-db-host",
	}

	var cfg integrationConfig
	if err := LoadWithFlags(path, &cfg, args); err != nil {
		t.Fatalf("LoadWithFlags returned error: %v", err)
	}

	if cfg.Name != "file-name" {
		t.Errorf("Name = %q, want file-name (only file sets it)", cfg.Name)
	}
	if cfg.Port != 4444 {
		t.Errorf("Port = %d, want 4444 (flag beats env and file)", cfg.Port)
	}
	if !reflect.DeepEqual(cfg.Tags, []string{"env-a", "env-b"}) {
		t.Errorf("Tags = %v, want [env-a env-b] (env beats file)", cfg.Tags)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s (default, nothing else sets it)", cfg.Timeout)
	}
	if cfg.DB.Host != "flag-db-host" {
		t.Errorf("DB.Host = %q, want flag-db-host (nested flag beats nested file)", cfg.DB.Host)
	}
	if cfg.DB.Port != 3333 {
		t.Errorf("DB.Port = %d, want 3333 (nested env beats nested file)", cfg.DB.Port)
	}
}

func TestFullPipelineJSONWithNestedAndDefaults(t *testing.T) {
	path := writeTempFile(t, "config.json", `{"name": "json-app", "db": {}}`)

	var cfg integrationConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Name != "json-app" {
		t.Errorf("Name = %q, want json-app", cfg.Name)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (default)", cfg.Port)
	}
	if cfg.DB.Host != "db-default" {
		t.Errorf("DB.Host = %q, want db-default (nested default)", cfg.DB.Host)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("DB.Port = %d, want 5432 (nested default, also satisfies nested required)", cfg.DB.Port)
	}
}
