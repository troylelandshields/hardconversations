package config

import (
	"io"

	yaml "gopkg.in/yaml.v3"
)

func parseConfig(rd io.Reader) (Config, error) {
	dec := yaml.NewDecoder(rd)
	dec.KnownFields(true)
	var conf Config
	if err := dec.Decode(&conf); err != nil {
		return conf, err
	}
	if conf.Version == "" {
		return conf, ErrMissingVersion
	}
	if conf.Version != "1" {
		return conf, ErrUnknownVersion
	}
	if len(conf.Conversations) == 0 {
		return conf, ErrNoPackages
	}
	for i := range conf.Conversations {
		for j := range conf.Conversations[i].Questions {
			inParsedType, err := conf.Conversations[i].Questions[j].Input.Parse()
			if err != nil {
				return conf, err
			}
			conf.Conversations[i].Questions[j].InputParsed = inParsedType

			outParsedType, err := conf.Conversations[i].Questions[j].Output.Parse()
			if err != nil {
				return conf, err
			}
			conf.Conversations[i].Questions[j].OutputParsed = outParsedType
		}
	}

	// if conf.Gen.Go != nil {
	// 	for i := range conf.Gen.Go.Overrides {
	// 		if err := conf.Gen.Go.Overrides[i].Parse(); err != nil {
	// 			return conf, err
	// 		}
	// 	}
	// }
	// for j := range conf.SQL {
	// 	if conf.SQL[j].Engine == "" {
	// 		return conf, ErrMissingEngine
	// 	}
	// 	if conf.SQL[j].Gen.Go != nil {
	// 		if conf.SQL[j].Gen.Go.Out == "" {
	// 			return conf, ErrNoPackagePath
	// 		}
	// 		if conf.SQL[j].Gen.Go.Package == "" {
	// 			conf.SQL[j].Gen.Go.Package = filepath.Base(conf.SQL[j].Gen.Go.Out)
	// 		}
	// 		for i := range conf.SQL[j].Gen.Go.Overrides {
	// 			if err := conf.SQL[j].Gen.Go.Overrides[i].Parse(); err != nil {
	// 				return conf, err
	// 			}
	// 		}
	// 	}
	// 	if conf.SQL[j].Gen.JSON != nil {
	// 		if conf.SQL[j].Gen.JSON.Out == "" {
	// 			return conf, ErrNoOutPath
	// 		}
	// 	}
	// 	for _, cg := range conf.SQL[j].Codegen {
	// 		if cg.Plugin == "" {
	// 			return conf, ErrPluginNoName
	// 		}
	// 		if cg.Out == "" {
	// 			return conf, ErrNoOutPath
	// 		}
	// 		// TOOD: Allow the use of built-in codegen from here
	// 		if _, ok := plugins[cg.Plugin]; !ok {
	// 			return conf, ErrPluginNotFound
	// 		}
	// 	}
	// }
	return conf, nil
}

// func (c *Config) validateGlobalOverrides() error {
// 	engines := map[Engine]struct{}{}
// 	for _, pkg := range c.SQL {
// 		if _, ok := engines[pkg.Engine]; !ok {
// 			engines[pkg.Engine] = struct{}{}
// 		}
// 	}
// 	if c.Gen.Go == nil {
// 		return nil
// 	}
// 	usesMultipleEngines := len(engines) > 1
// 	for _, oride := range c.Gen.Go.Overrides {
// 		if usesMultipleEngines && oride.Engine == "" {
// 			return fmt.Errorf(`the "engine" field is required for global type overrides because your configuration uses multiple database engines`)
// 		}
// 	}
// 	return nil
// }