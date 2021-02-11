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

	return applyEnvOverridesToStruct(v)
}

func applyEnvOverridesToStruct(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			if err := applyEnvOverridesToStruct(fieldVal); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		}

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		if err := setFieldFromString(fieldVal, envVal); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}

	return nil
}
