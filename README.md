# go-config-loader

Loads configuration from a YAML or JSON file into a struct, then overlays matching environment variables on top (env vars win).

## Usage

```go
type AppConfig struct {
	Port  int    `yaml:"port" env:"PORT" default:"8080"`
	Host  string `yaml:"host" env:"HOST" validate:"required"`
	Debug bool   `yaml:"debug" env:"DEBUG"`
}

var cfg AppConfig
if err := config.Load("config.yaml", &cfg); err != nil {
	log.Fatal(err)
}
```

`config.Load` reads `config.yaml` (or `.json`), unmarshals it into `cfg`, checks each field's `env` tag and overrides it if that environment variable is set, then fills any still-zero field from its `default` tag, and finally checks `validate:"required"` tags. See the sections below for the full tag reference, the command-line flag overlay, nested structs, and `Watch`.

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

## Watching for changes

`Watch` loads the file once, then watches it for changes (using [fsnotify](https://github.com/fsnotify/fsnotify)) and reloads it in place whenever it changes:

```go
var cfg AppConfig
stop, err := config.Watch("config.yaml", &cfg, os.Args[1:], func(err error) {
	if err != nil {
		log.Printf("reload failed: %v", err)
		return
	}
	log.Printf("config reloaded: %+v", cfg)
})
if err != nil {
	log.Fatal(err)
}
defer stop()
```

Every reload runs the same pipeline as `LoadWithFlags` — **env var and command-line flag overlays are reapplied on every reload**, not just the initial load, so a flag passed at startup or an env var set in the process's environment continues to take precedence over whatever the file now says. Pass `nil` for `args` to skip flag parsing, matching `Load`'s behavior.

Watch debounces rapid bursts of filesystem events (common when editors save via a temp-file-plus-rename) into a single reload, and watches the file's containing directory so renames are picked up correctly. Call the returned `stop` function to stop watching and release the underlying OS resources; after `stop` returns, `onChange` will not be called again (though a reload already in progress may still finish).

`out` is mutated in place on every reload — if your program reads `cfg` from multiple goroutines, guard access with a mutex or copy it inside `onChange` before use.

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

## Running the demo

`main.go` loads `config.yaml` into a nested `AppConfig` (which embeds a `DBConfig`), using `default`, `env`, and `flag` tags throughout, then intentionally loads a second, incomplete config to show the `*config.MissingFieldsError` output:

```sh
go run .
```

```
Loaded config: {Port:8080 Host:localhost Debug:false Tags:[web prod] Timeout:30s DB:{Host:db.internal Port:5432}}

Try: PORT=9090 go run .                 — env var override
Try: go run . -port=9191                — flag override (highest precedence)
Try: go run . -db-host=db.example.com    — override a nested field

--- Demonstrating a validation error ---
Got expected error: missing required fields: [Host Port]
```

Try the suggested overrides, including one on the nested `DB.Host` field via its `flag:"db-host"` tag:

```sh
go run . -port=9191 -db-host=db.example.com
PORT=9090 go run .
```
