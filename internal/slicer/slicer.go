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
// and treats the command as failed if it returns a non-zero exit code OR if
// it produces any output on the stderr stream.
func (e CmdExecutor) Run(workDir, command string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	cmd.Stderr = &stderr

	err := cmd.Run()

	// rsync can exit with code 0 but still report errors (e.g., permission denied)
	// on stderr. We must treat any stderr output as a failure.
	if err == nil && stderr.Len() > 0 {
		return fmt.Errorf("command produced stderr output:\n%s", stderr.String())
	}
	if err != nil {
		return fmt.Errorf("command failed with exit code: %w\nSTDERR:\n%s", err, stderr.String())
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
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temporary file %s: %v\n", tmpFile.Name(), err)
		}
	}()

	if _, err := tmpFile.WriteString(strings.Join(files, "\n")); err != nil {
		return fmt.Errorf("failed to write to temporary manifest file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary manifest file: %w", err)
	}

	args := []string{
		"-a",
		"--include-from=" + tmpFile.Name(),
		"--include=*/",
		"--exclude=*",
		".",
		output,
	}

	if err := exec.Run(source, "rsync", args...); err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}

	return nil
}
