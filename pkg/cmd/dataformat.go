package cmd

import (
	"errors"
	"strings"
)

// A readDataFormat is a format that chezmoi uses for reading (JSON, TOML, or
// YAML) and implements the github.com/spf13/pflag.Value interface.
type readDataFormat string

// A writeDataFormat is format that chezmoi uses for writing (JSON or YAML) and
// implements the github.com/spf13/pflag.Value interface.
//
// TOML is not included as write format as it requires the top level value to be
// an object, and chezmoi occasionally requires the top level value to be a
// simple value or array.
type writeDataFormat string

const (
	readDataFormatJSON readDataFormat = "json"
	readDataFormatTOML readDataFormat = "toml"
	readDataFormatYAML readDataFormat = "yaml"

	writeDataFormatJSON writeDataFormat = "json"
	writeDataFormatYAML writeDataFormat = "yaml"

	defaultWriteDataFormat = writeDataFormatJSON
)

// Set implements github.com/spf13/pflag.Value.Set.
func (f *readDataFormat) Set(s string) error {
	switch strings.ToLower(s) {
	case "json":
		*f = readDataFormatJSON
	case "toml":
		*f = readDataFormatTOML
	case "yaml":
		*f = readDataFormatYAML
	default:
		return errors.New("invalid or unsupported data format")
	}
	return nil
}

// String implements github.com/spf13/pflag.Value.String.
func (f readDataFormat) String() string {
	return string(f)
}

// Type implements github.com/spf13/pflag.Value.Type.
func (f readDataFormat) Type() string {
	return "json|toml|yaml"
}

// Set implements github.com/spf13/pflag.Value.Set.
func (f *writeDataFormat) Set(s string) error {
	switch strings.ToLower(s) {
	case "json":
		*f = writeDataFormatJSON
	case "yaml":
		*f = writeDataFormatYAML
	default:
		return errors.New("invalid or unsupported data format")
	}
	return nil
}

// String implements github.com/spf13/pflag.Value.String.
func (f writeDataFormat) String() string {
	return string(f)
}

// Type implements github.com/spf13/pflag.Value.Type.
func (f writeDataFormat) Type() string {
	return "json|yaml"
}
