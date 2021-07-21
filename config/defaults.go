package config

import (
	"fmt"
	"reflect"
)

// applyDefaults walks the fields of out's underlying struct and, for each
// field tagged with `default:"..."`, sets the field to that value if the
// field is still at its zero value. It recurses into nested structs.
func applyDefaults(out interface{}) error {
	v, err := structValue(out)
	if err != nil {
		return err
	}

	return walkFields(v, "", func(field reflect.StructField, fieldVal reflect.Value, path string) error {
		defaultVal, ok := field.Tag.Lookup("default")
		if !ok {
			return nil
		}

		if !fieldVal.IsZero() {
			return nil
		}

		if err := setFieldFromString(fieldVal, defaultVal); err != nil {
			return fmt.Errorf("field %s: %w", path, err)
		}

		return nil
	})
}
