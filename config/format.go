package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Format identifies the serialization format of a config file.
type Format string

const (
	// FormatYAML parses the file as YAML.
	FormatYAML Format = "yaml"
	// FormatJSON parses the file as JSON.
	FormatJSON Format = "json"
)

// detectFormat infers a Format from a file's extension. It returns an
// error for unrecognized extensions so callers can fail fast instead of
// silently misparsing a file.
func detectFormat(path string) (Format, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("cannot detect format from extension of %q", path)
	}
}
