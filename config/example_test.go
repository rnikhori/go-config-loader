package config_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rnikhori/go-config-loader/config"
)

type ExampleAppConfig struct {
	Port int    `yaml:"port" env:"EXAMPLE_PORT" default:"8080"`
	Host string `yaml:"host" env:"EXAMPLE_HOST" default:"localhost"`
}

func ExampleLoad() {
	dir, err := os.MkdirTemp("", "go-config-loader-example")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("host: example.com\n"), 0o644); err != nil {
		fmt.Println("error:", err)
		return
	}

	var cfg ExampleAppConfig
	if err := config.Load(path, &cfg); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Port=%d Host=%s\n", cfg.Port, cfg.Host)
	// Output: Port=8080 Host=example.com
}
