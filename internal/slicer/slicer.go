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

// Executor defines an interface for running external commands.
type Executor interface {
	Run(workDir, command string, args ...string) error
}

// TempFiler defines an interface for creating and managing temporary files.
type TempFiler interface {
	CreateTemp(dir, pattern string) (TempFile, error)
}

// TempFile defines an interface for interacting with a temporary file.
type TempFile interface {
	io.WriteCloser
	Name() string
}

// liveTempFiler is a concrete implementation of TempFiler using the os package.
type liveTempFiler struct{}

func (ltf *liveTempFiler) CreateTemp(dir, pattern string) (TempFile, error) {
	return os.CreateTemp(dir, pattern)
}

// CmdExecutor is a concrete implementation of the Executor interface.
type CmdExecutor struct{}

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
func Slice(source, output string, files []string, exec Executor, filer TempFiler) error {
	tmpFile, err := filer.CreateTemp("", "repo-slice-manifest-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary manifest file: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temporary file %s: %v\n", tmpFile.Name(), err)
		}
	}()

	if _, err := tmpFile.Write([]byte(strings.Join(files, "\n"))); err != nil {
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
