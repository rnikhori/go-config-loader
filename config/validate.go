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
	walkFields(v, "", func(field reflect.StructField, fieldVal reflect.Value, path string) error {
		if field.Tag.Get("validate") == "required" && fieldVal.IsZero() {
			missing = append(missing, path)
		}
		return nil
	})

	if len(missing) > 0 {
		return &MissingFieldsError{Fields: missing}
	}

	return nil
}
