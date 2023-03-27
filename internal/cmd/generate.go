package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/trace"
	"strings"
	"sync"

	"github.com/troylelandshields/hardconversations/internal/compiler"
	"github.com/troylelandshields/hardconversations/internal/config"
	"github.com/troylelandshields/hardconversations/internal/multierr"
	"golang.org/x/sync/errgroup"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 or 2 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The supported version can only be "1".
`

const errMessageNoPackages = `No packages are configured`

func printFileErr(stderr io.Writer, dir string, fileErr *multierr.FileError) {
	filename := strings.TrimPrefix(fileErr.Filename, dir+"/")
	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
}

func readConfig(stderr io.Writer, dir, filename string) (string, *config.Config, error) {
	configPath := ""
	if filename != "" {
		configPath = filepath.Join(dir, filename)
	} else {
		var yamlMissing, jsonMissing bool
		yamlPath := filepath.Join(dir, "hardc.yaml")
		jsonPath := filepath.Join(dir, "hardc.json")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			yamlMissing = true
		}
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			jsonMissing = true
		}

		if yamlMissing && jsonMissing {
			fmt.Fprintln(stderr, "error parsing configuration files. sqlc.yaml or sqlc.json: file does not exist")
			return "", nil, errors.New("config file missing")
		}

		if !yamlMissing && !jsonMissing {
			fmt.Fprintln(stderr, "error: both sqlc.json and sqlc.yaml files present")
			return "", nil, errors.New("sqlc.json and sqlc.yaml present")
		}

		configPath = yamlPath
		if yamlMissing {
			configPath = jsonPath
		}
	}

	base := filepath.Base(configPath)
	blob, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}

	conf, err := config.ParseConfig(bytes.NewReader(blob))
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprintf(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprintf(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprintf(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing %s: %s\n", base, err)
		return "", nil, err
	}

	return configPath, &conf, nil
}

func Generate(ctx context.Context, dir, filename string, stderr io.Writer) (map[string]string, error) {
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(configPath)
	if err := config.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return nil, err
	}

	output := map[string]string{}
	errored := false

	var m sync.Mutex
	grp, gctx := errgroup.WithContext(ctx)
	grp.SetLimit(runtime.GOMAXPROCS(0))

	stderrs := make([]bytes.Buffer, len(conf.Conversations))

	for i, lConvo := range conf.Conversations {
		convo := lConvo
		errout := &stderrs[i]

		grp.Go(func() error {
			// combo := config.Combine(*conf, sql.SQL)
			// if sql.Plugin != nil {
			// 	combo.Codegen = *sql.Plugin
			// }

			// // TODO: This feels like a hack that will bite us later
			// joined := make([]string, 0, len(sql.Schema))
			// for _, s := range sql.Schema {
			// 	joined = append(joined, filepath.Join(dir, s))
			// }
			// sql.Schema = joined

			// joined = make([]string, 0, len(sql.Queries))
			// for _, q := range sql.Queries {
			// 	joined = append(joined, filepath.Join(dir, q))
			// }
			// sql.Queries = joined

			packageRegion := trace.StartRegion(gctx, "package")
			trace.Logf(gctx, "", "dir=%s", dir)

			compiler := parse(gctx, dir, convo, errout)

			resp, err := compiler.Compile(gctx)
			if err != nil {
				fmt.Fprintf(errout, "# \n")
				fmt.Fprintf(errout, "error generating code: %s\n", err)
				errored = true
				packageRegion.End()
				return nil
			}

			files := map[string]string{}
			for _, file := range resp {
				files[file.Name] = string(file.Contents)
			}

			m.Lock()
			for n, source := range files {
				filename := filepath.Join(dir, n)
				output[filename] = source
			}
			m.Unlock()

			packageRegion.End()
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return nil, err
	}
	if errored {
		for i, _ := range stderrs {
			if _, err := io.Copy(stderr, &stderrs[i]); err != nil {
				return nil, err
			}
		}
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}

func parse(ctx context.Context, dir string, convo config.Conversation, stderr io.Writer) *compiler.Compiler {
	return compiler.NewCompiler(convo)
}
