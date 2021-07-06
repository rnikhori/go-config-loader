package config

import (
	"errors"
	"testing"
)

type testRequiredConfig struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required"`
}

type testRequiredNestedConfig struct {
	Name string             `yaml:"name" validate:"required"`
	DB   testRequiredDBPart `yaml:"db"`
}

type testRequiredDBPart struct {
	Host string `yaml:"host" validate:"required"`
}

func TestLoadRequiredFieldMissing(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: localhost\n")

	var cfg testRequiredConfig
	err := Load(path, &cfg)
	if err == nil {
		t.Fatal("expected error for missing required field, got nil")
	}

	var missingErr *MissingFieldsError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected *MissingFieldsError, got %T: %v", err, err)
	}
	if len(missingErr.Fields) != 1 || missingErr.Fields[0] != "Port" {
		t.Errorf("Fields = %v, want [Port]", missingErr.Fields)
	}
}

func TestLoadRequiredFieldsAllSatisfied(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: localhost\nport: 8080\n")

	var cfg testRequiredConfig
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
}

func TestLoadRequiredListsAllMissingFields(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "{}\n")

	var cfg testRequiredConfig
	err := Load(path, &cfg)

	var missingErr *MissingFieldsError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected *MissingFieldsError, got %T: %v", err, err)
	}
	if len(missingErr.Fields) != 2 {
		t.Errorf("Fields = %v, want 2 entries", missingErr.Fields)
	}
}

func TestLoadRequiredNestedFieldMissing(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "name: myapp\ndb:\n  host: \"\"\n")

	var cfg testRequiredNestedConfig
	err := Load(path, &cfg)

	var missingErr *MissingFieldsError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected *MissingFieldsError, got %T: %v", err, err)
	}
	if len(missingErr.Fields) != 1 || missingErr.Fields[0] != "DB.Host" {
		t.Errorf("Fields = %v, want [DB.Host]", missingErr.Fields)
	}
}

func TestLoadRequiredSatisfiedByEnv(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: \"\"\nport: 8080\n")

	type cfgType struct {
		Host string `yaml:"host" env:"TEST_REQUIRED_HOST" validate:"required"`
		Port int    `yaml:"port" validate:"required"`
	}

	t.Setenv("TEST_REQUIRED_HOST", "envhost")

	var cfg cfgType
	if err := Load(path, &cfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
}

func TestLoadRequiredSatisfiedByDefault(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: \"\"\n")

	type cfgType struct {
		Host string `yaml:"host" validate:"required"`
		Port int    `yaml:"port" default:"8080" validate:"required"`
	}

	var cfg cfgType
	err := Load(path, &cfg)

	var missingErr *MissingFieldsError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected *MissingFieldsError, got %T: %v", err, err)
	}
	if len(missingErr.Fields) != 1 || missingErr.Fields[0] != "Host" {
		t.Errorf("Fields = %v, want [Host] (Port satisfied by default)", missingErr.Fields)
	}
}
