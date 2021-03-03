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

	return applyDefaultsToStruct(v)
}

func applyDefaultsToStruct(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			if err := applyDefaultsToStruct(fieldVal); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		}

		defaultVal, ok := field.Tag.Lookup("default")
		if !ok {
			continue
		}

		if !fieldVal.IsZero() {
			continue
		}

		if err := setFieldFromString(fieldVal, defaultVal); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}

	return nil
}
