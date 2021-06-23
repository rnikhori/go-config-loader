package config

import (
	"fmt"
	"reflect"
	"strings"
)

// MissingFieldsError is returned by Load (and its variants) when one or
// more fields tagged `validate:"required"` are still at their zero value
// after all overlay layers have been applied.
type MissingFieldsError struct {
	Fields []string
}

func (e *MissingFieldsError) Error() string {
	return fmt.Sprintf("config: missing required fields: %s", strings.Join(e.Fields, ", "))
}

// validateRequired walks out's underlying struct (recursing into nested
// structs) and collects the dotted path of every field tagged
// `validate:"required"` that is still at its zero value. If any are
// found, it returns a *MissingFieldsError.
func validateRequired(out interface{}) error {
	v, err := structValue(out)
	if err != nil {
		return err
	}

	var missing []string
	collectMissingRequired(v, "", &missing)

	if len(missing) > 0 {
		return &MissingFieldsError{Fields: missing}
	}

	return nil
}

func collectMissingRequired(v reflect.Value, prefix string, missing *[]string) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)
		path := field.Name
		if prefix != "" {
			path = prefix + "." + path
		}

		if fieldVal.Kind() == reflect.Struct {
			collectMissingRequired(fieldVal, path, missing)
		}

		if field.Tag.Get("validate") == "required" && fieldVal.IsZero() {
			*missing = append(*missing, path)
		}
	}
}
