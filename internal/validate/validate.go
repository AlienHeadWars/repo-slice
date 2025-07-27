// file: internal/validate/validate.go

// Package validate provides functions for validating the input configuration
// for the repo-slice tool.
package validate

import (
	"fmt"
	"io/fs"
	"os"
)

// FS defines an interface for file system operations, allowing for mock
// implementations in tests. This adheres to the Dependency Inversion Principle.
type FS interface {
	Stat(name string) (fs.FileInfo, error)
}

// LiveFS is a concrete implementation of the FS interface that uses the
// standard library's os package.
type LiveFS struct{}

// Stat performs the Stat operation using the os package.
func (l LiveFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// Config represents the input configuration that needs validation.
type Config struct {
	ManifestPath string
	SourcePath   string
}

// ValidateInputs checks that the source and manifest paths are valid.
// It depends on the FS interface, not a concrete file system.
func ValidateInputs(cfg Config, fsys FS) error {
	// Validate Source Path
	sourceInfo, err := fsys.Stat(cfg.SourcePath)
	if err != nil {
		return fmt.Errorf("source path '%s' not found", cfg.SourcePath)
	}
	if !sourceInfo.IsDir() {
		return fmt.Errorf("source path '%s' is not a directory", cfg.SourcePath)
	}

	// Validate Manifest Path
	manifestInfo, err := fsys.Stat(cfg.ManifestPath)
	if err != nil {
		return fmt.Errorf("manifest path '%s' not found", cfg.ManifestPath)
	}
	if manifestInfo.IsDir() {
		return fmt.Errorf("manifest path '%s' is a directory, not a file", cfg.ManifestPath)
	}

	return nil
}
