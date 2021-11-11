package config

import "testing"

type testFlagConfig struct {
	Port int    `yaml:"port" env:"TEST_FLAG_PORT" flag:"port"`
	Host string `yaml:"host" env:"TEST_FLAG_HOST" flag:"host"`
}

type testFlagNestedConfig struct {
	Name string           `yaml:"name" flag:"name"`
	DB   testFlagDBNested `yaml:"db"`
}

type testFlagDBNested struct {
	Host string `yaml:"host" flag:"db-host"`
}

func TestLoadWithFlagsOverridesEnv(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: filehost\n")

	t.Setenv("TEST_FLAG_PORT", "9090")

	var cfg testFlagConfig
	if err := LoadWithFlags(path, &cfg, []string{"-port=7070"}); err != nil {
		t.Fatalf("LoadWithFlags returned error: %v", err)
	}

	if cfg.Port != 7070 {
		t.Errorf("Port = %d, want 7070 (flag beats env and file)", cfg.Port)
	}
	if cfg.Host != "filehost" {
		t.Errorf("Host = %q, want filehost (no flag passed, file value kept)", cfg.Host)
	}
}

func TestLoadWithFlagsOnlyAppliesPassedFlags(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: filehost\n")

	var cfg testFlagConfig
	if err := LoadWithFlags(path, &cfg, []string{"-host=cliHost"}); err != nil {
		t.Fatalf("LoadWithFlags returned error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (no -port flag passed)", cfg.Port)
	}
	if cfg.Host != "cliHost" {
		t.Errorf("Host = %q, want cliHost", cfg.Host)
	}
}

func TestLoadWithFlagsNestedStruct(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "name: fromfile\ndb:\n  host: filehost\n")

	var cfg testFlagNestedConfig
	if err := LoadWithFlags(path, &cfg, []string{"-db-host=clihost"}); err != nil {
		t.Fatalf("LoadWithFlags returned error: %v", err)
	}

	if cfg.Name != "fromfile" {
		t.Errorf("Name = %q, want fromfile", cfg.Name)
	}
	if cfg.DB.Host != "clihost" {
		t.Errorf("DB.Host = %q, want clihost", cfg.DB.Host)
	}
}

func TestLoadWithFlagsNilArgsBehavesLikeLoad(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: filehost\n")

	var cfg testFlagConfig
	if err := LoadWithFlags(path, &cfg, nil); err != nil {
		t.Fatalf("LoadWithFlags returned error: %v", err)
	}

	if cfg.Port != 8080 || cfg.Host != "filehost" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadFormatWithFlagsJSON(t *testing.T) {
	path := writeTempFile(t, "config.cfg", `{"port": 8080, "host": "jsonhost"}`)

	var cfg testFlagConfig
	if err := LoadFormatWithFlags(path, FormatJSON, &cfg, []string{"-port=5555"}); err != nil {
		t.Fatalf("LoadFormatWithFlags returned error: %v", err)
	}

	if cfg.Port != 5555 {
		t.Errorf("Port = %d, want 5555 (flag beats JSON file value)", cfg.Port)
	}
	if cfg.Host != "jsonhost" {
		t.Errorf("Host = %q, want jsonhost (no flag passed for host)", cfg.Host)
	}
}

type testDuplicateFlagConfig struct {
	A string `yaml:"a" flag:"same"`
	B string `yaml:"b" flag:"same"`
}

func TestApplyFlagOverridesDuplicateFlagNameErrors(t *testing.T) {
	var cfg testDuplicateFlagConfig
	err := applyFlagOverrides(&cfg, []string{"-same=x"})
	if err == nil {
		t.Fatal("expected error for duplicate flag name, got nil")
	}
}

type testBoolFlagConfig struct {
	Debug bool `yaml:"debug" flag:"debug"`
}

// TestBoolFlagRequiresExplicitValue documents a real constraint: because
// every flag is registered as a string flag (so setFieldFromString can
// convert it to the field's actual type later), a bool field cannot be
// toggled with a standalone -debug the way flag.Bool would allow. It
// must be passed as -debug=true.
func TestBoolFlagRequiresExplicitValue(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "debug: false\n")

	var cfg testBoolFlagConfig
	if err := LoadWithFlags(path, &cfg, []string{"-debug"}); err == nil {
		t.Fatal("expected error for standalone -debug (no value), got nil")
	}
}

func TestBoolFlagWithExplicitValue(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "debug: false\n")

	var cfg testBoolFlagConfig
	if err := LoadWithFlags(path, &cfg, []string{"-debug=true"}); err != nil {
		t.Fatalf("LoadWithFlags returned error: %v", err)
	}
	if !cfg.Debug {
		t.Error("Debug = false, want true")
	}
}

func TestLoadWithFlagsUnknownFlagErrors(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\n")

	var cfg testFlagConfig
	if err := LoadWithFlags(path, &cfg, []string{"-nope=1"}); err == nil {
		t.Fatal("expected error for unregistered flag, got nil")
	}
}
