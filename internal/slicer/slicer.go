// file: internal/slicer/slicer.go
package slicer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	if err == nil && stderr.Len() > 0 {
		return fmt.Errorf("command produced stderr output:\n%s", stderr.String())
	}
	if err != nil {
		return fmt.Errorf("command failed with exit code: %w\nSTDERR:\n%s", err, stderr.String())
	}
	return nil
}

// ParseManifest reads and parses a manifest from any io.Reader.
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
func Slice(source, output, manifestPath string, exec Executor) error {
	// Implementation to be added in a future step.
	// This stub ensures the code compiles with the new signature.
	return nil
}