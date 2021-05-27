package config

import (
	"flag"
	"fmt"
	"io"
	"reflect"
)

// applyFlagOverrides registers a flag for every field tagged with
// `flag:"name"` (recursing into nested structs), parses args against
// them, and sets only the fields whose flags were actually passed on the
// command line. Fields without a corresponding flag on the command line
// are left untouched.
func applyFlagOverrides(out interface{}, args []string) error {
	v, err := structValue(out)
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	bindings := map[string]reflect.Value{}
	registerFlags(fs, v, bindings)

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	var applyErr error
	fs.Visit(func(f *flag.Flag) {
		if applyErr != nil {
			return
		}
		fieldVal, ok := bindings[f.Name]
		if !ok {
			return
		}
		if err := setFieldFromString(fieldVal, f.Value.String()); err != nil {
			applyErr = fmt.Errorf("flag %s: %w", f.Name, err)
		}
	})

	return applyErr
}

// registerFlags walks v's fields, registering a string flag for each one
// tagged with `flag:"name"` and recording its name-to-field binding.
// Nested struct fields are recursed into. Every flag is registered as a
// string flag regardless of the underlying field type; the raw value is
// converted by setFieldFromString once we know which flags were actually
// set by the caller.
func registerFlags(fs *flag.FlagSet, v reflect.Value, bindings map[string]reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			registerFlags(fs, fieldVal, bindings)
		}

		flagName := field.Tag.Get("flag")
		if flagName == "" {
			continue
		}

		bindings[flagName] = fieldVal
		fs.String(flagName, "", fmt.Sprintf("overrides %s", field.Name))
	}
}
