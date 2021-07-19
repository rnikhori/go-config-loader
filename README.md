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

## Command-line flags

A field tagged with `flag:"name"` can be overridden from the command line, taking precedence over env vars, the file, and defaults:

```go
type AppConfig struct {
	Port int `yaml:"port" env:"PORT" flag:"port" default:"8080"`
}

var cfg AppConfig
config.LoadWithFlags("config.yaml", &cfg, os.Args[1:])
```

```sh
go run . -port=9090
```

Only flags actually present in the given args are applied — an absent flag leaves the field as set by the env/file/default layers. `Load` and `LoadFormat` do not parse flags at all (equivalent to passing nil args); use `LoadWithFlags` or `LoadFormatWithFlags` to opt in.

Full precedence, highest to lowest: **flags > env > file > defaults**.

## Nested structs

Nested (including embedded) struct fields are walked recursively, so each level of nesting can carry its own tags:

```go
type DBConfig struct {
	Host string `yaml:"host" env:"DB_HOST"`
	Port int    `yaml:"port" env:"DB_PORT" default:"5432"`
}

type AppConfig struct {
	Name string   `yaml:"name" env:"APP_NAME"`
	DB   DBConfig `yaml:"db"`
}
```

## Defaults

A field tagged with `default:"..."` is filled in with that value if it is still at its zero value after the file and env layers have been applied:

```go
type AppConfig struct {
	Port int `yaml:"port" env:"PORT" default:"8080"`
}
```

Since the default only applies when a field is at its zero value, a default cannot be distinguished from an explicit zero in the file (e.g. `port: 0`) — both result in the default being applied. Use a non-zero sentinel or a pointer type if that distinction matters for your field.

## Required fields

A field tagged with `validate:"required"` must end up non-zero after every other layer (file, env, flags, defaults) has been applied. If any required fields are still unset, `Load` returns a `*config.MissingFieldsError` listing all of them (not just the first one found), with nested fields reported as dotted paths like `DB.Host`:

```go
type AppConfig struct {
	Host string `yaml:"host" validate:"required"`
}

var cfg AppConfig
if err := config.Load("config.yaml", &cfg); err != nil {
	var missing *config.MissingFieldsError
	if errors.As(err, &missing) {
		log.Fatalf("missing required config: %v", missing.Fields)
	}
	log.Fatal(err)
}
```

A default satisfies a required field as long as it leaves the field non-zero, since validation runs after defaults are applied.

## Supported field kinds

Besides strings, bools, and numeric types, the following kinds are supported for file parsing, env overrides, and defaults:

| Go type         | File syntax (YAML/JSON) | Env/default syntax         |
|-----------------|--------------------------|------------------------------|
| `[]string`      | native list              | comma-separated, e.g. `"a,b,c"` |
| `time.Duration` | string, e.g. `"30s"`     | string, e.g. `"30s"` (parsed with `time.ParseDuration`) |

```go
type AppConfig struct {
	Tags    []string      `yaml:"tags" env:"TAGS"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" default:"30s"`
}
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
