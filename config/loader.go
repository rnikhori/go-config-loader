// Package config loads configuration from a file, overlays matching
// environment variables on top, and fills any still-unset fields with
// their declared defaults. Precedence, highest to lowest: env vars,
// file values, defaults.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads the config file at path, detecting its format (YAML or JSON)
// from the file extension, into out, then overrides any field tagged with
// `env:"NAME"` using the value of that environment variable, if set. out
// must be a non-nil pointer to a struct.
func Load(path string, out interface{}) error {
	format, err := detectFormat(path)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	return LoadFormat(path, format, out)
}

// LoadFormat behaves like Load but uses the given format explicitly
// instead of inferring it from the file extension. This is useful when a
// config file does not have a conventional extension.
func LoadFormat(path string, format Format, out interface{}) error {
	if _, err := structValue(out); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config: read file: %w", err)
	}

	switch format {
	case FormatYAML:
		if err := yaml.Unmarshal(data, out); err != nil {
			return fmt.Errorf("config: parse yaml: %w", err)
		}
	case FormatJSON:
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("config: parse json: %w", err)
		}
	default:
		return fmt.Errorf("config: unsupported format %q", format)
	}

	if err := applyEnvOverrides(out); err != nil {
		return fmt.Errorf("config: apply env overrides: %w", err)
	}

	if err := applyDefaults(out); err != nil {
		return fmt.Errorf("config: apply defaults: %w", err)
	}

	return nil
}
