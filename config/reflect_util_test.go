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
