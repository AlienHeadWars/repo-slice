// file: internal/slicer/slicer.go
package slicer

import (
	"bytes"
	"fmt"
	"os/exec"
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

// Slice constructs and executes an rsync command to copy files based on a
// manifest file that uses rsync filter-rule syntax.
func Slice(source, output, manifestPath string, exec Executor) error {
	args := []string{
		"-a", // Archive mode to preserve permissions, ownership, etc.
		"--filter",
		fmt.Sprintf("merge %s", manifestPath),
		".",    // Source directory (relative to the workDir)
		output, // Destination directory
	}

	if err := exec.Run(source, "rsync", args...); err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}

	return nil
}
