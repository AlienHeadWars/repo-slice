// file: internal/slicer/slicer.go
package slicer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Executor defines an interface for running external commands from a specific
// working directory.
type Executor interface {
	Run(workDir, command string, args ...string) error
}

// CmdExecutor is a concrete implementation of the Executor interface.
type CmdExecutor struct{}

// Run executes a command from the given working directory. It captures stderr
// and includes it in the error message if the command fails.
func (e CmdExecutor) Run(workDir, command string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir // Set the command's working directory.
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed: %w\nSTDERR:\n%s", err, stderr.String())
	}
	return nil
}

// ParseManifest reads and parses a manifest from any io.Reader. It returns a
// slice of file paths, ignoring empty lines and lines starting with '#'.
func ParseManifest(r io.Reader) ([]string, error) {
	var paths []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
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

	// The `-aR` combination (archive and relative) is the most robust way to
	// copy only the specified files while preserving their directory structure.
	// By executing from within the source directory, we ensure the paths in
	// the manifest are interpreted correctly.
	args := []string{
		"-aR",
		"--files-from=" + tmpFile.Name(),
		".", // The source is the current directory (which we set to `source`).
		output,
	}

	if err := exec.Run(source, "rsync", args...); err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}

	return nil
}
