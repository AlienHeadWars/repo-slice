// file: internal/slicer/slicer.go

// Package slicer contains the core logic for parsing the manifest and
// executing the rsync command to create a repository slice.
package slicer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

// Executor defines an interface for running external commands. This allows for
// a mock implementation to be used in unit tests, adhering to the Dependency
// Inversion Principle.
type Executor interface {
	// Run executes a command with the given arguments and returns an error
	// if the command fails.
	Run(command string, args ...string) error
}

// CmdExecutor is a concrete implementation of the Executor interface that runs
// real commands on the operating system.
type CmdExecutor struct{}

// Run executes a command using the os/exec package.
func (e CmdExecutor) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ParseManifest reads and parses a manifest file from the provided file system.
// It returns a slice of file paths, ignoring empty lines and trimming whitespace.
func (fsys *validate.LiveFS) ParseManifest(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open manifest file: %w", err)
	}
	defer file.Close()

	return parseManifestReader(file)
}

// parseManifestReader is a helper that parses content from any io.Reader.
func parseManifestReader(r io.Reader) ([]string, error) {
	var paths []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			paths = append(paths, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading manifest content: %w", err)
	}

	return paths, nil
}

// Slice constructs and executes an rsync command to copy files from a source
// to a destination directory, using a temporary file created from the manifest list.
func Slice(source, output string, files []string, exec Executor) error {
	// Create a temporary file to hold the list of files for rsync.
	tmpFile, err := os.CreateTemp("", "repo-slice-manifest-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary manifest file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the file paths to the temporary file.
	if _, err := tmpFile.WriteString(strings.Join(files, "\n")); err != nil {
		return fmt.Errorf("failed to write to temporary manifest file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary manifest file: %w", err)
	}

	// Construct and run the rsync command.
	args := []string{
		"-a", // Archive mode: recursive, preserves symlinks, permissions, times, etc.
		"--files-from=" + tmpFile.Name(),
		source,
		output,
	}

	if err := exec.Run("rsync", args...); err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}

	return nil
}
