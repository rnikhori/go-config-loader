package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rnikhori/go-config-loader/config"
)

type DBConfig struct {
	Host string `yaml:"host" json:"host" env:"DB_HOST" flag:"db-host" default:"localhost"`
	Port int    `yaml:"port" json:"port" env:"DB_PORT" flag:"db-port" default:"5432"`
}

type AppConfig struct {
	Port    int           `yaml:"port" json:"port" env:"PORT" flag:"port" default:"8080"`
	Host    string        `yaml:"host" json:"host" env:"HOST" flag:"host" validate:"required"`
	Debug   bool          `yaml:"debug" json:"debug" env:"DEBUG" flag:"debug"`
	Tags    []string      `yaml:"tags" json:"tags" env:"TAGS" flag:"tags"`
	Timeout time.Duration `yaml:"timeout" json:"timeout" env:"TIMEOUT" flag:"timeout" default:"30s"`
	DB      DBConfig      `yaml:"db" json:"db"`
}

func main() {
	var cfg AppConfig
	if err := config.LoadWithFlags("config.yaml", &cfg, os.Args[1:]); err != nil {
		exitOnLoadError(err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)
	fmt.Println()
	fmt.Println("Try: PORT=9090 go run .                 — env var override")
	fmt.Println("Try: go run . -port=9191                — flag override (highest precedence)")
	fmt.Println("Try: go run . -db-host=db.example.com    — override a nested field")
	fmt.Println()

	demoValidationError()
}

// exitOnLoadError prints a *config.MissingFieldsError specially, since
// it carries a structured list of every missing required field rather
// than just a single message.
func exitOnLoadError(err error) {
	var missing *config.MissingFieldsError
	if errors.As(err, &missing) {
		log.Fatalf("config is missing required fields: %v", missing.Fields)
	}
	log.Fatalf("failed to load config: %v", err)
}

// demoValidationError shows what happens when a validate:"required" field
// is left unset: Load returns a *config.MissingFieldsError listing every
// missing field, not just the first one found.
func demoValidationError() {
	fmt.Println("--- Demonstrating a validation error ---")

	type Incomplete struct {
		Host string `yaml:"host" validate:"required"`
		Port int    `yaml:"port" validate:"required"`
	}

	dir, err := os.MkdirTemp("", "go-config-loader-demo")
	if err != nil {
		log.Printf("demo: could not create temp dir: %v", err)
		return
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "incomplete.yaml")
	if err := os.WriteFile(path, []byte("host: \"\"\n"), 0o644); err != nil {
		log.Printf("demo: could not write temp file: %v", err)
		return
	}

	var incomplete Incomplete
	err = config.Load(path, &incomplete)

	var missing *config.MissingFieldsError
	if errors.As(err, &missing) {
		fmt.Printf("Got expected error: missing required fields: %v\n", missing.Fields)
		return
	}

	fmt.Printf("unexpected result loading intentionally incomplete config: %v\n", err)
}
