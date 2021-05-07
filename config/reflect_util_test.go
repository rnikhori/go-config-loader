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
