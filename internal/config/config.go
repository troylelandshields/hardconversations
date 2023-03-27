package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

type versionSetting struct {
	Number string `json:"version" yaml:"version"`
}

type Engine string

type Paths []string

func (p *Paths) UnmarshalJSON(data []byte) error {
	if string(data[0]) == `[` {
		var out []string
		if err := json.Unmarshal(data, &out); err != nil {
			return nil
		}
		*p = Paths(out)
		return nil
	}
	var out string
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	*p = Paths([]string{out})
	return nil
}

func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out := []string{}
	if sliceErr := unmarshal(&out); sliceErr != nil {
		var ele string
		if strErr := unmarshal(&ele); strErr != nil {
			return strErr
		}
		out = []string{ele}
	}

	*p = Paths(out)
	return nil
}

type Config struct {
	Version       string         `json:"version" yaml:"version"`
	Conversations []Conversation `json:"conversations" yaml:"conversations"`
}

type Conversation struct {
	Path        string     `json:"path" yaml:"path"`
	Instruction string     `json:"instruction" yaml:"instruction"`
	Questions   []Question `json:"questions" yaml:"questions"`
}

type Question struct {
	FunctionName string `json:"function_name" yaml:"function_name"`
	Prompt       string `json:"prompt" yaml:"prompt"`
	Input        GoType `json:"input" yaml:"input"`
	Output       GoType `json:"output" yaml:"output"`

	InputParsed  *ParsedGoType
	OutputParsed *ParsedGoType
}

var ErrMissingEngine = errors.New("unknown engine")
var ErrMissingVersion = errors.New("no version number")
var ErrNoOutPath = errors.New("no output path")
var ErrNoPackageName = errors.New("missing package name")
var ErrNoPackagePath = errors.New("missing package path")
var ErrNoPackages = errors.New("no packages")
var ErrNoQuerierType = errors.New("no querier emit type enabled")
var ErrUnknownEngine = errors.New("invalid engine")
var ErrUnknownVersion = errors.New("invalid version number")

var ErrPluginBuiltin = errors.New("a built-in plugin with that name already exists")
var ErrPluginNoName = errors.New("missing plugin name")
var ErrPluginExists = errors.New("a plugin with that name already exists")
var ErrPluginNotFound = errors.New("no plugin found")
var ErrPluginNoType = errors.New("plugin: field `process` or `wasm` required")
var ErrPluginBothTypes = errors.New("plugin: both `process` and `wasm` cannot both be defined")
var ErrPluginProcessNoCmd = errors.New("plugin: missing process command")

var ErrInvalidQueryParameterLimit = errors.New("invalid query parameter limit")

func ParseConfig(rd io.Reader) (Config, error) {
	var buf bytes.Buffer
	var config Config
	var version versionSetting

	ver := io.TeeReader(rd, &buf)
	dec := yaml.NewDecoder(ver)
	if err := dec.Decode(&version); err != nil {
		return config, err
	}
	if version.Number == "" {
		return config, ErrMissingVersion
	}
	switch version.Number {
	case "1":
		return parseConfig(&buf)
	default:
		return config, ErrUnknownVersion
	}
}

func Validate(c *Config) error {
	return nil
}

// type CombinedSettings struct {
// 	Global    Config
// 	Package   SQL
// 	Go        SQLGo
// 	JSON      SQLJSON
// 	Rename    map[string]string
// 	Overrides []Override

// 	// TODO: Combine these into a more usable type
// 	Codegen Codegen
// }

// func Combine(conf Config, pkg SQL) CombinedSettings {
// 	cs := CombinedSettings{
// 		Global:  conf,
// 		Package: pkg,
// 	}
// 	if conf.Gen.Go != nil {
// 		cs.Rename = conf.Gen.Go.Rename
// 		cs.Overrides = append(cs.Overrides, conf.Gen.Go.Overrides...)
// 	}
// 	if pkg.Gen.Go != nil {
// 		cs.Go = *pkg.Gen.Go
// 		cs.Overrides = append(cs.Overrides, pkg.Gen.Go.Overrides...)
// 	}
// 	if pkg.Gen.JSON != nil {
// 		cs.JSON = *pkg.Gen.JSON
// 	}
// 	return cs
// }
