// file: internal/slicer/slicer.go
package slicer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Executor defines an interface for running external commands.
type Executor interface {
	// Run executes a command from a specific working directory.
	Run(workDir, command string, args ...string) error
}

// CmdExecutor is a concrete implementation of the Executor interface.
type CmdExecutor struct{}

// Run executes a command from the given working directory.
func (e CmdExecutor) Run(workDir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ParseManifest reads and parses a manifest from any io.Reader.
func ParseManifest(r io.Reader) ([]string, error) {
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

// Slice constructs and executes an rsync command to copy files.
func Slice(source, output string, files []string, exec Executor) error {
	tmpFile, err := os.CreateTemp("", "repo-slice-manifest-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary manifest file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(strings.Join(files, "\n")); err != nil {
		return fmt.Errorf("failed to write to temporary manifest file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary manifest file: %w", err)
	}

	// The `-R` (relative) flag is crucial. It tells rsync to preserve the
	// path structure of the files listed in the manifest.
	args := []string{
		"-a",
		"-R",
		"--files-from=" + tmpFile.Name(),
		".", // Source is the current directory (which we set to `source`).
		output,
	}

	// Execute the command from within the source directory.
	if err := exec.Run(source, "rsync", args...); err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}

	return nil
}
