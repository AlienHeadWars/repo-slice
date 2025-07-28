// file: cmd/repo-slice/main_test.go
package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlienHeadwars/repo-slice/internal/slicer"
)

const (
	testDirSource    = "valid-source"
	testFileManifest = "valid-manifest.txt"
	testFileSource   = "source-is-a-file"
	testFileOutput   = "test-output"

	flagSource   = "--source"
	flagManifest = "--manifest"
	flagOutput   = "--output"
)

// mockExecutor implements the slicer.Executor interface for testing.
type mockExecutor struct {
	returnErr bool
}

func (m *mockExecutor) Run(workDir, command string, args ...string) error {
	if m.returnErr {
		return errors.New("mock executor error")
	}
	// For happy-path tests, use the real executor to ensure integration.
	realExecutor := slicer.CmdExecutor{}
	return realExecutor.Run(workDir, command, args...)
}

// setupTestFS creates a temporary directory structure for testing.
func setupTestFS(t *testing.T) string {
	t.Helper()
	rootDir, err := os.MkdirTemp("", "repo-slice-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(rootDir); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	})

	sourcePath := filepath.Join(rootDir, testDirSource)
	if err := os.Mkdir(sourcePath, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourcePath, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatalf("failed to create a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourcePath, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatalf("failed to create b.txt: %v", err)
	}

	manifestContent := "a.txt\n"
	if err := os.WriteFile(filepath.Join(rootDir, testFileManifest), []byte(manifestContent), 0644); err != nil {
		t.Fatalf("failed to create manifest file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(rootDir, testFileSource), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	return rootDir
}

// TestRunEndToEnd verifies the full application logic.
func TestRunEndToEnd(t *testing.T) {
	rootDir := setupTestFS(t)
	outputPath := filepath.Join(rootDir, testFileOutput)

	t.Run("Valid run creates correct output", func(t *testing.T) {
		args := []string{
			flagSource, filepath.Join(rootDir, testDirSource),
			flagManifest, filepath.Join(rootDir, testFileManifest),
			flagOutput, outputPath,
		}
		executor := &mockExecutor{returnErr: false}

		err := run(args, executor)
		if err != nil {
			t.Fatalf("run() with valid args failed unexpectedly: %v", err)
		}

		if _, err := os.Stat(filepath.Join(outputPath, "a.txt")); os.IsNotExist(err) {
			t.Error("expected file 'a.txt' was not found in the output directory")
		}
		if _, err := os.Stat(filepath.Join(outputPath, "b.txt")); !os.IsNotExist(err) {
			t.Error("unexpected file 'b.txt' was found in the output directory")
		}
	})

	t.Run("Slice operation fails", func(t *testing.T) {
		args := []string{
			flagSource, filepath.Join(rootDir, testDirSource),
			flagManifest, filepath.Join(rootDir, testFileManifest),
			flagOutput, outputPath,
		}
		// Configure the mock to return an error.
		executor := &mockExecutor{returnErr: true}

		err := run(args, executor)
		if err == nil {
			t.Error("run() with a failing slicer did not return an error")
		}
	})
}
