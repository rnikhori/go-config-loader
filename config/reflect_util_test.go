package config

import (
	"reflect"
	"testing"
	"time"
)

func TestSetFieldFromStringDuration(t *testing.T) {
	var d time.Duration
	v := reflect.ValueOf(&d).Elem()

	if err := setFieldFromString(v, "30s"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 30*time.Second {
		t.Errorf("duration = %v, want 30s", d)
	}
}

func TestSetFieldFromStringDurationInvalid(t *testing.T) {
	var d time.Duration
	v := reflect.ValueOf(&d).Elem()

	if err := setFieldFromString(v, "not-a-duration"); err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
}

func TestSetFieldFromStringStringSlice(t *testing.T) {
	var tags []string
	v := reflect.ValueOf(&tags).Elem()

	if err := setFieldFromString(v, "a, b ,c"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("tags = %v, want %v", tags, want)
	}
}

func TestSetFieldFromStringStringSliceEmpty(t *testing.T) {
	var tags []string
	v := reflect.ValueOf(&tags).Elem()

	if err := setFieldFromString(v, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("tags = %v, want empty slice", tags)
	}
}

func TestSetFieldFromStringUnsupportedSliceElement(t *testing.T) {
	var nums []int
	v := reflect.ValueOf(&nums).Elem()

	if err := setFieldFromString(v, "1,2,3"); err == nil {
		t.Fatal("expected error for unsupported slice element type, got nil")
	}
}

func TestSetFieldFromStringScalarKinds(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		var s string
		if err := setFieldFromString(reflect.ValueOf(&s).Elem(), "hello"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s != "hello" {
			t.Errorf("s = %q, want hello", s)
		}
	})

	t.Run("int8", func(t *testing.T) {
		var n int8
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "42"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 42 {
			t.Errorf("n = %d, want 42", n)
		}
	})

	t.Run("int16", func(t *testing.T) {
		var n int16
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "-100"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != -100 {
			t.Errorf("n = %d, want -100", n)
		}
	})

	t.Run("int32", func(t *testing.T) {
		var n int32
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "70000"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 70000 {
			t.Errorf("n = %d, want 70000", n)
		}
	})

	t.Run("int64", func(t *testing.T) {
		var n int64
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "9000000000"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 9000000000 {
			t.Errorf("n = %d, want 9000000000", n)
		}
	})

	t.Run("int invalid", func(t *testing.T) {
		var n int
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "not-a-number"); err == nil {
			t.Fatal("expected error for invalid int, got nil")
		}
	})

	t.Run("uint8", func(t *testing.T) {
		var n uint8
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "200"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 200 {
			t.Errorf("n = %d, want 200", n)
		}
	})

	t.Run("uint invalid", func(t *testing.T) {
		var n uint
		if err := setFieldFromString(reflect.ValueOf(&n).Elem(), "-1"); err == nil {
			t.Fatal("expected error for negative uint, got nil")
		}
	})

	t.Run("bool true", func(t *testing.T) {
		var b bool
		if err := setFieldFromString(reflect.ValueOf(&b).Elem(), "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !b {
			t.Error("b = false, want true")
		}
	})

	t.Run("bool invalid", func(t *testing.T) {
		var b bool
		if err := setFieldFromString(reflect.ValueOf(&b).Elem(), "yes"); err == nil {
			t.Fatal("expected error for invalid bool, got nil")
		}
	})

	t.Run("float32", func(t *testing.T) {
		var f float32
		if err := setFieldFromString(reflect.ValueOf(&f).Elem(), "3.5"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f != 3.5 {
			t.Errorf("f = %v, want 3.5", f)
		}
	})

	t.Run("float64 invalid", func(t *testing.T) {
		var f float64
		if err := setFieldFromString(reflect.ValueOf(&f).Elem(), "not-a-float"); err == nil {
			t.Fatal("expected error for invalid float, got nil")
		}
	})

	t.Run("unsupported kind", func(t *testing.T) {
		var m map[string]string
		if err := setFieldFromString(reflect.ValueOf(&m).Elem(), "x"); err == nil {
			t.Fatal("expected error for unsupported kind (map), got nil")
		}
	})
}

func TestSetFieldFromStringUnsettable(t *testing.T) {
	type s struct {
		unexported string
	}
	var v s
	rv := reflect.ValueOf(&v).Elem().Field(0)

	if err := setFieldFromString(rv, "x"); err == nil {
		t.Fatal("expected error for unexported (unsettable) field, got nil")
	}
}
