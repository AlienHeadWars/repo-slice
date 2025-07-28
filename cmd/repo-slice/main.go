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

	"github.com/AlienHeadwars/repo-slice/internal/slicer"
	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

// ... (Config struct remains the same) ...
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
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	fsys := &validate.LiveFS{}
	validationCfg := validate.Config{
		SourcePath:   cfg.SourcePath,
		ManifestPath: cfg.ManifestPath,
	}

	if err := validate.ValidateInputs(validationCfg, fsys); err != nil {
		return err
	}

	// Open the manifest file to get an io.Reader.
	manifestFile, err := os.Open(cfg.ManifestPath)
	if err != nil {
		return fmt.Errorf("could not open manifest file: %w", err)
	}
	defer manifestFile.Close()

	// Parse the manifest file.
	files, err := slicer.ParseManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Execute the slice operation.
	executor := slicer.CmdExecutor{}
	if err := slicer.Slice(cfg.SourcePath, cfg.OutputPath, files, executor); err != nil {
		return fmt.Errorf("failed to execute slice operation: %w", err)
	}

	fmt.Printf("Successfully created repository slice in %s\n", cfg.OutputPath)
	return nil
}

// ... (parseArgs function remains the same) ...
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
