package config

import (
	"os"
	"testing"
	"time"
)

type testWatchConfig struct {
	Port int    `yaml:"port" env:"TEST_WATCH_PORT"`
	Host string `yaml:"host"`
}

func TestWatchReloadsOnFileChange(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: initial\n")

	var cfg testWatchConfig
	changes := make(chan error, 10)

	stop, err := Watch(path, &cfg, nil, func(err error) {
		changes <- err
	})
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	defer stop()

	if cfg.Host != "initial" {
		t.Fatalf("initial load: Host = %q, want initial", cfg.Host)
	}

	if err := os.WriteFile(path, []byte("port: 9090\nhost: updated\n"), 0o644); err != nil {
		t.Fatalf("write updated file: %v", err)
	}

	select {
	case err := <-changes:
		if err != nil {
			t.Fatalf("onChange called with error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for reload callback")
	}

	if cfg.Port != 9090 || cfg.Host != "updated" {
		t.Fatalf("after reload: unexpected config: %+v", cfg)
	}
}

func TestWatchEnvOverlayAppliesAfterReload(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: initial\n")

	t.Setenv("TEST_WATCH_PORT", "7777")

	var cfg testWatchConfig
	changes := make(chan error, 10)

	stop, err := Watch(path, &cfg, nil, func(err error) {
		changes <- err
	})
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	defer stop()

	if cfg.Port != 7777 {
		t.Fatalf("initial load: Port = %d, want 7777 (env override)", cfg.Port)
	}

	if err := os.WriteFile(path, []byte("port: 1\nhost: updated\n"), 0o644); err != nil {
		t.Fatalf("write updated file: %v", err)
	}

	select {
	case err := <-changes:
		if err != nil {
			t.Fatalf("onChange called with error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for reload callback")
	}

	if cfg.Port != 7777 {
		t.Fatalf("after reload: Port = %d, want 7777 (env override still applied)", cfg.Port)
	}
	if cfg.Host != "updated" {
		t.Fatalf("after reload: Host = %q, want updated", cfg.Host)
	}
}

func TestWatchStopPreventsFurtherCallbacks(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "port: 8080\nhost: initial\n")

	var cfg testWatchConfig
	changes := make(chan error, 10)

	stop, err := Watch(path, &cfg, nil, func(err error) {
		changes <- err
	})
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	if err := stop(); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}

	if err := os.WriteFile(path, []byte("port: 2\nhost: afterstop\n"), 0o644); err != nil {
		t.Fatalf("write updated file: %v", err)
	}

	select {
	case err := <-changes:
		t.Fatalf("onChange called after stop: %v", err)
	case <-time.After(300 * time.Millisecond):
	}
}

func TestWatchInvalidPathErrors(t *testing.T) {
	var cfg testWatchConfig
	if _, err := Watch("does-not-exist.yaml", &cfg, nil, func(error) {}); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
