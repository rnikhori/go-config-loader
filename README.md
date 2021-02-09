# go-config-loader

Loads configuration from a YAML or JSON file into a struct, then overlays matching environment variables on top (env vars win).

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

## File formats

`Load` detects the file format from its extension:

| Extension        | Format |
|-------------------|--------|
| `.yaml`, `.yml`    | YAML   |
| `.json`            | JSON   |

```go
var cfg AppConfig
config.Load("config.json", &cfg) // parsed as JSON
```

A struct intended to be loaded from either format should carry both `yaml` and `json` tags where the field names differ from the default (JSON matches field names case-insensitively when no `json` tag is present, so simple configs work with only `yaml` tags).

If a file's extension isn't `.yaml`, `.yml`, or `.json`, use `config.LoadFormat` to specify the format explicitly:

```go
config.LoadFormat("config.cfg", config.FormatJSON, &cfg)
```

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
