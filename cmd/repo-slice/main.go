// file: cmd/repo-slice/main.go

/*
Repo-slice is a command-line tool that automates the creation of
streamlined, context-specific branches from a larger repository.

It reads a manifest file (an "allow-list") and creates a clean,
filtered copy of a source directory, which can then be pushed to a
dedicated branch for use by AI assistants.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

// Config holds the configuration options for the repo-slice tool,
// parsed from command-line arguments.
type Config struct {
	ManifestPath string
	SourcePath   string
	OutputPath   string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run executes the main logic of the application based on the provided arguments.
func run(args []string) error {
	parsedCfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	fsys := validate.LiveFS{}
	validationCfg := validate.Config{
		SourcePath:   parsedCfg.SourcePath,
		ManifestPath: parsedCfg.ManifestPath,
	}

	if err := validate.ValidateInputs(validationCfg, fsys); err != nil {
		return err
	}

	return nil
}

// parseArgs parses the command-line arguments from a string slice and returns
// a populated Config struct. It returns an error if the arguments cannot be
// parsed.
func parseArgs(args []string) (Config, error) {
	var cfg Config
	fs := flag.NewFlagSet("repo-slice", flag.ContinueOnError)

	fs.StringVar(&cfg.ManifestPath, "manifest", "", "Path to the manifest file (required)")
	fs.StringVar(&cfg.SourcePath, "source", ".", "The source directory to read from")
	fs.StringVar(&cfg.OutputPath, "output", "", "The destination directory (required)")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
