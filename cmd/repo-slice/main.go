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
	"io"
	"os"

	"github.com/AlienHeadwars/repo-slice/internal/slicer"
	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

// Config holds the configuration options for the repo-slice tool.
type Config struct {
	ManifestPath string
	SourcePath   string
	OutputPath   string
}

// FileSystem defines an interface for file system operations needed by run.
type FileSystem interface {
	ValidateInputs(cfg validate.Config) error
	Open(name string) (io.ReadCloser, error)
}

// Slicer defines an interface for the core application logic.
type Slicer interface {
	ParseManifest(r io.Reader) ([]string, error)
	Slice(source, output string, files []string) error
}

// liveFS is a concrete implementation of the FileSystem interface.
type liveFS struct{}

func (fs *liveFS) ValidateInputs(cfg validate.Config) error {
	return validate.ValidateInputs(cfg, &validate.LiveFS{})
}
func (fs *liveFS) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// liveSlicer is a concrete implementation of the Slicer interface.
type liveSlicer struct{}

func (s *liveSlicer) ParseManifest(r io.Reader) ([]string, error) {
	return slicer.ParseManifest(r)
}
func (s *liveSlicer) Slice(source, output string, files []string) error {
	// Inject the concrete dependencies here.
	executor := &slicer.CmdExecutor{}
	filer := &slicer.LiveTempFiler{}
	return slicer.Slice(source, output, files, executor, filer)
}

func main() {
	if err := run(os.Args[1:], &liveFS{}, &liveSlicer{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run executes the main logic of the application.
func run(args []string, fsys FileSystem, slicer Slicer) (err error) {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	validationCfg := validate.Config{
		SourcePath:   cfg.SourcePath,
		ManifestPath: cfg.ManifestPath,
	}
	if err := fsys.ValidateInputs(validationCfg); err != nil {
		return err
	}

	manifestFile, err := fsys.Open(cfg.ManifestPath)
	if err != nil {
		return fmt.Errorf("could not open manifest file: %w", err)
	}
	defer func() {
		if closeErr := manifestFile.Close(); err == nil {
			err = closeErr
		}
	}()

	files, err := slicer.ParseManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	if err := slicer.Slice(cfg.SourcePath, cfg.OutputPath, files); err != nil {
		return fmt.Errorf("failed to execute slice operation: %w", err)
	}

	fmt.Printf("Successfully created repository slice in %s\n", cfg.OutputPath)
	return nil
}

// parseArgs parses the command-line arguments.
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
