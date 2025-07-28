//go:build integration

package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCmdExecutor_RunFails is an integration test for the CmdExecutor that
// verifies its error handling for both invalid flags and permission errors.
func TestCmdExecutor_RunFails(t *testing.T) {
	executor := CmdExecutor{}

	t.Run("invalid flag", func(t *testing.T) {
		err := executor.Run(".", "rsync", "--non-existent-flag")
		if err == nil {
			t.Fatal("CmdExecutor.Run() did not return an error for a failing command")
		}
		if !strings.Contains(err.Error(), "STDERR") {
			t.Error("Error message from failed command did not include STDERR")
		}
	})

	t.Run("permission denied", func(t *testing.T) {
		unwritableDir, err := os.MkdirTemp("", "unwritable-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(unwritableDir)
		if err := os.Chmod(unwritableDir, 0555); err != nil {
			t.Fatalf("Failed to chmod temp dir: %v", err)
		}

		err = executor.Run(".", "touch", filepath.Join(unwritableDir, "test.txt"))
		if err == nil {
			t.Fatal("CmdExecutor.Run() did not return an error for a permission denied scenario")
		}
	})
}

func TestSliceIntegration(t *testing.T) {
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

	_ = os.WriteFile(filepath.Join(sourceDir, "a.txt"), []byte("a"), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "b.txt"), []byte("b"), 0644)
	_ = os.Mkdir(filepath.Join(sourceDir, "subdir"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "subdir", "c.txt"), []byte("c"), 0644)

	filesToCopy := []string{"a.txt", "subdir/c.txt"}

	// Use the real CmdExecutor for the integration test.
	executor := CmdExecutor{}
	err = Slice(sourceDir, outputDir, filesToCopy, executor)
	if err != nil {
		t.Fatalf("Slice() failed during integration test: %v", err)
	}

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
