package config

import (
	"fmt"
	"os"
)

// applyEnvOverrides walks the fields of out's underlying struct and, for
// each field tagged with `env:"NAME"`, sets the field from the named
// environment variable if it is present.
func applyEnvOverrides(out interface{}) error {
	v, err := structValue(out)
	if err != nil {
		return err
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		if err := setFieldFromString(v.Field(i), envVal); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}

	return nil
}
