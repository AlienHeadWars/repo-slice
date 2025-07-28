//go:build integration

// This file contains integration tests that interact with the real file system
// and execute the rsync command. To run these tests, use the build tag
// 'integration': go test -v ./... -tags=integration

package slicer

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// liveExecutor is a real implementation of the Executor interface.
type liveExecutor struct{}

func (e liveExecutor) Run(command string, args ...string) error {
	return exec.Command(command, args...).Run()
}

func TestSliceIntegration(t *testing.T) {
	// Setup: Create temporary source directory, files, and manifest
	sourceDir, err := os.MkdirTemp("", "source-*")
	if err != nil {
		t.Fatalf("Failed to create temp source dir: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	outputDir, err := os.MkdirTemp("", "output-*")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Create files in the source directory
	_ = os.WriteFile(filepath.Join(sourceDir, "a.txt"), []byte("a"), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "b.txt"), []byte("b"), 0644)
	_ = os.Mkdir(filepath.Join(sourceDir, "subdir"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "subdir", "c.txt"), []byte("c"), 0644)

	// Create the manifest content
	manifestContent := "a.txt\nsubdir/c.txt\n"
	filesToCopy := []string{"a.txt", "subdir/c.txt"}

	// Execution: Run the slice operation with a real executor
	executor := liveExecutor{}
	err = Slice(sourceDir, outputDir, filesToCopy, executor)
	if err != nil {
		t.Fatalf("Slice() failed during integration test: %v", err)
	}

	// Assertions: Check the contents of the output directory
	if _, err := os.Stat(filepath.Join(outputDir, "a.txt")); os.IsNotExist(err) {
		t.Error("Expected file 'a.txt' was not found in the output directory")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "subdir", "c.txt")); os.IsNotExist(err) {
		t.Error("Expected file 'subdir/c.txt' was not found in the output directory")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "b.txt")); !os.IsNotExist(err) {
		t.Error("Unexpected file 'b.txt' was found in the output directory")
	}
}