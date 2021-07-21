package config

import (
	"fmt"
	"os"
	"reflect"
)

// applyEnvOverrides walks the fields of out's underlying struct and, for
// each field tagged with `env:"NAME"`, sets the field from the named
// environment variable if it is present. Nested struct fields (including
// embedded ones) are recursed into so their fields can carry their own
// env tags.
func applyEnvOverrides(out interface{}) error {
	v, err := structValue(out)
	if err != nil {
		return err
	}

	return walkFields(v, "", func(field reflect.StructField, fieldVal reflect.Value, path string) error {
		envKey := field.Tag.Get("env")
		if envKey == "" {
			return nil
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			return nil
		}

		if err := setFieldFromString(fieldVal, envVal); err != nil {
			return fmt.Errorf("field %s: %w", path, err)
		}

		return nil
	})
}
