package main

import (
	"fmt"
	"log"

	"github.com/rnikhori/go-config-loader/config"
)

type AppConfig struct {
	Port  int    `yaml:"port" env:"PORT"`
	Host  string `yaml:"host" env:"HOST"`
	Debug bool   `yaml:"debug" env:"DEBUG"`
}

func main() {
	var cfg AppConfig
	if err := config.Load("config.yaml", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)
	fmt.Println("Try: PORT=9090 go run . — to see env vars override the YAML value.")
}
