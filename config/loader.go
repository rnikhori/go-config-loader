// Package config loads configuration from a file and overlays
// matching environment variables on top, with env vars taking precedence.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads the YAML file at path into out, then overrides any field
// tagged with `env:"NAME"` using the value of that environment variable,
// if set. out must be a non-nil pointer to a struct.
func Load(path string, out interface{}) error {
	if _, err := structValue(out); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config: read file: %w", err)
	}

	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := applyEnvOverrides(out); err != nil {
		return fmt.Errorf("config: apply env overrides: %w", err)
	}

	return nil
}
