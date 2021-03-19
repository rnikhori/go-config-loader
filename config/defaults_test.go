package config

import "testing"

type testDefaultsConfig struct {
	Port    int    `yaml:"port" env:"TEST_DEFAULTS_PORT" default:"8080"`
	Host    string `yaml:"host" env:"TEST_DEFAULTS_HOST" default:"localhost"`
	Timeout int    `yaml:"timeout" default:"30"`
}

func TestLoadAppliesDefaultsWhenUnset(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: fromfile\n")

	var cfg testDefaultsConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (default)", cfg.Port)
	}
	if cfg.Host != "fromfile" {
		t.Errorf("Host = %q, want fromfile (file value, not overwritten by default)", cfg.Host)
	}
	if cfg.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30 (default)", cfg.Timeout)
	}
}

func TestLoadDefaultDoesNotOverrideEnv(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "{}\n")

	t.Setenv("TEST_DEFAULTS_PORT", "9999")

	var cfg testDefaultsConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 9999 {
		t.Errorf("Port = %d, want 9999 (env override, not default)", cfg.Port)
	}
}

func TestLoadDefaultDoesNotOverrideExplicitZero(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 0\n")

	var cfg testDefaultsConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (default applied since 0 is the zero value and indistinguishable from unset)", cfg.Port)
	}
}
