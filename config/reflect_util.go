package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var durationType = reflect.TypeOf(time.Duration(0))

// structValue validates that out is a non-nil pointer to a struct and
// returns the dereferenced struct value.
func structValue(out interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return reflect.Value{}, fmt.Errorf("out must be a non-nil pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("out must point to a struct")
	}

	return v, nil
}

// setFieldFromString converts raw into the appropriate type for field and
// sets it. field must be addressable and settable.
func setFieldFromString(field reflect.Value, raw string) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	if field.Type() == durationType {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("parse duration: %w", err)
		}
		field.SetInt(int64(d))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Slice:
		return setSliceFromString(field, raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int: %w", err)
		}
		field.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("parse uint: %w", err)
		}
		field.SetUint(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("parse bool: %w", err)
		}
		field.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return fmt.Errorf("parse float: %w", err)
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("unsupported field type %s", field.Kind())
	}

	return nil
}

// setSliceFromString parses raw as a comma-separated list and sets field
// to the resulting slice. Only []string is currently supported.
func setSliceFromString(field reflect.Value, raw string) error {
	if field.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("unsupported slice element type %s", field.Type().Elem().Kind())
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		return nil
	}

	parts := strings.Split(raw, ",")
	values := make([]string, len(parts))
	for i, p := range parts {
		values[i] = strings.TrimSpace(p)
	}

	slice := reflect.MakeSlice(field.Type(), len(values), len(values))
	for i, s := range values {
		slice.Index(i).SetString(s)
	}
	field.Set(slice)

	return nil
}
