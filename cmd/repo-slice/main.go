// file: cmd/repo-slice/main.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

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

	// The `main` function is the composition root. Here we "inject" the
	// real file system implementation into our validation logic.
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

// parseArgs parses the command-line arguments and returns a Config struct.
func parseArgs(args []string) (Config, error) {
	var cfg Config
	fs := flag.NewFlagSet("repo-slice", flag.ContinueOnError)

	fs.StringVar(&cfg.ManifestPath, "manifest", "", "Path to the manifest file (required)")
	fs.StringVar(&cfg.SourcePath, "source", ".", "The source directory to read from")
	fs.StringVar(&cfg.OutputPath, "output", "", "The destination directory (required)")

	// This call to Parse will handle parsing the arguments and printing usage
	// information if -h or -help is provided.
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	return cfg, nil
}