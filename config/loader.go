// Package config loads configuration from a file, overlays matching
// environment variables and command-line flags on top, fills any
// still-unset fields with their declared defaults, and validates that
// required fields ended up set. Precedence, highest to lowest: flags,
// env vars, file values, defaults.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads the config file at path, detecting its format (YAML or JSON)
// from the file extension, into out, then overrides any field tagged
// with `env:"NAME"` using the value of that environment variable, and
// finally fills any field still at its zero value with its `default:"..."`
// tag, if present. out must be a non-nil pointer to a struct.
//
// Load does not apply command-line flag overrides; use LoadWithFlags for
// that.
func Load(path string, out interface{}) error {
	return LoadWithFlags(path, out, nil)
}

// LoadFormat behaves like Load but uses the given format explicitly
// instead of inferring it from the file extension.
func LoadFormat(path string, format Format, out interface{}) error {
	return LoadFormatWithFlags(path, format, out, nil)
}

// LoadWithFlags behaves like Load, additionally overlaying command-line
// flags parsed from args on top of the env layer. args is typically
// os.Args[1:]. A field is only overridden if its flag was actually
// present in args; flags are registered from `flag:"name"` tags. A nil
// args is equivalent to calling Load.
func LoadWithFlags(path string, out interface{}, args []string) error {
	format, err := detectFormat(path)
	if err != nil {
		return wrapStage("detect format", err)
	}

	return LoadFormatWithFlags(path, format, out, args)
}

// LoadFormatWithFlags combines the explicit-format behavior of LoadFormat
// with the flag overlay behavior of LoadWithFlags.
func LoadFormatWithFlags(path string, format Format, out interface{}, args []string) error {
	if _, err := structValue(out); err != nil {
		return wrapStage("validate target", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return wrapStage("read file", err)
	}

	switch format {
	case FormatYAML:
		if err := yaml.Unmarshal(data, out); err != nil {
			return wrapStage("parse yaml", err)
		}
	case FormatJSON:
		if err := json.Unmarshal(data, out); err != nil {
			return wrapStage("parse json", err)
		}
	default:
		return fmt.Errorf("config: unsupported format %q", format)
	}

	if err := applyEnvOverrides(out); err != nil {
		return wrapStage("apply env overrides", err)
	}

	if args != nil {
		if err := applyFlagOverrides(out, args); err != nil {
			return wrapStage("apply flag overrides", err)
		}
	}

	if err := applyDefaults(out); err != nil {
		return wrapStage("apply defaults", err)
	}

	if err := validateRequired(out); err != nil {
		return err
	}

	return nil
}
