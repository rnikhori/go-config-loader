package config

import (
	"fmt"
	"os"
	"reflect"
)

// applyEnvOverrides walks the fields of out's underlying struct and, for
// each field tagged with `env:"NAME"`, sets the field from the named
// environment variable if it is present.
func applyEnvOverrides(out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("out must be a non-nil pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("out must point to a struct")
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
