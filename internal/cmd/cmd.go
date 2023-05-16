package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/trace"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	// "github.com/troylelandshields/hardc/internal/codegen/golang"
	// "github.com/troylelandshields/hardc/internal/debug"
	// "github.com/troylelandshields/hardc/internal/info"
	// "github.com/troylelandshields/hardc/internal/tracer"
)

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "hardc", SilenceUsage: true}
	rootCmd.PersistentFlags().StringP("file", "f", "", "specify an alternate config file (default: hardc.yaml)")

	rootCmd.AddCommand(genCmd)
	// rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SilenceErrors = true

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return 1
		}
	}
	return 0
}

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the sqlc version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "version").End()
		if version == "" {
			fmt.Printf("%s\n", "0.0.7")
		} else {
			fmt.Printf("%s\n", version)
		}
		return nil
	},
}

//	var initCmd = &cobra.Command{
//		Use:   "init",
//		Short: "Create an empty hardc.yaml settings file",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			defer trace.StartRegion(cmd.Context(), "init").End()
//			file := "hardc.yaml"
//			if f := cmd.Flag("file"); f != nil && f.Changed {
//				file = f.Value.String()
//				if file == "" {
//					return fmt.Errorf("file argument is empty")
//				}
//			}
//			if _, err := os.Stat(file); !os.IsNotExist(err) {
//				return nil
//			}
//			blob, err := yaml.Marshal(config.V1GenerateSettings{Version: "1"})
//			if err != nil {
//				return err
//			}
//			return os.WriteFile(file, blob, 0644)
//		},
//	}

func getConfigPath(stderr io.Writer, f *pflag.Flag) (string, string) {
	if f != nil && f.Changed {
		file := f.Value.String()
		if file == "" {
			fmt.Fprintln(stderr, "error parsing config: file argument is empty")
			os.Exit(1)
		}
		abspath, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintf(stderr, "error parsing config: absolute file path lookup failed: %s\n", err)
			os.Exit(1)
		}
		return filepath.Dir(abspath), filepath.Base(abspath)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
			os.Exit(1)
		}
		return wd, ""
	}
}

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code from SQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		output, err := Generate(cmd.Context(), dir, name, stderr)
		if err != nil {
			return err
		}
		for filename, source := range output {
			os.MkdirAll(filepath.Dir(filename), 0755)
			if err := os.WriteFile(filename, []byte(source), 0644); err != nil {
				fmt.Fprintf(stderr, "%s: %s\n", filename, err)
				return err
			}
		}
		return nil
	},
}

func getLines(f []byte) []string {
	fp := bytes.NewReader(f)
	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
