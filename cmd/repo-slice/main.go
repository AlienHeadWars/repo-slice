// file: cmd/repo-slice/main.go
package main

import (
	"flag"
	"fmt"
	"os"
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
	// The config is currently unused but will be validated and used
	// in subsequent features.
	_, err := parseArgs(args)
	if err != nil {
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