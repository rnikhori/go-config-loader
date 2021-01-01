# go-config-loader

Loads configuration from a YAML file into a struct, then overlays matching environment variables on top (env vars win).

## Usage

```go
type AppConfig struct {
	Port  int    `yaml:"port" env:"PORT"`
	Host  string `yaml:"host" env:"HOST"`
	Debug bool   `yaml:"debug" env:"DEBUG"`
}

var cfg AppConfig
if err := config.Load("config.yaml", &cfg); err != nil {
	log.Fatal(err)
}
```

`config.Load` reads `config.yaml`, unmarshals it into `cfg`, then checks each field's `env` tag and, if that environment variable is set, overrides the field's value.

## Running

```sh
go run .
```

With `config.yaml` containing `port: 8080`, `host: localhost`, `debug: false`, this prints:

```
Loaded config: {Port:8080 Host:localhost Debug:false}
```

Override a value via environment variable:

```sh
PORT=9090 DEBUG=true go run .
```

```
Loaded config: {Port:9090 Host:localhost Debug:true}
```
