// file: cmd/repo-slice/main_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
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

// setupTestFS creates a temporary directory structure for testing. It registers
// a cleanup function with the test to automatically remove the directory
// when the test completes. It returns the root path of the created directory.
func setupTestFS(t *testing.T) string {
	t.Helper()
	rootDir, err := os.MkdirTemp("", "repo-slice-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// t.Cleanup registers a function to be called when the test
	// and all its subtests complete. This is the idiomatic way
	// to handle test cleanup.
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

// TestRunEndToEnd is an end-to-end test for the main run function. It verifies
// the full application logic, from argument parsing and validation to the
// final slice operation.
func TestRunEndToEnd(t *testing.T) {
	rootDir := setupTestFS(t)
	outputPath := filepath.Join(rootDir, testFileOutput)

	t.Run("Valid run creates correct output", func(t *testing.T) {
		args := []string{
			flagSource, filepath.Join(rootDir, testDirSource),
			flagManifest, filepath.Join(rootDir, testFileManifest),
			flagOutput, outputPath,
		}

		err := run(args)
		if err != nil {
			t.Fatalf("run() with valid args failed unexpectedly: %v", err)
		}

		// Verify the output directory
		if _, err := os.Stat(filepath.Join(outputPath, "a.txt")); os.IsNotExist(err) {
			t.Error("expected file 'a.txt' was not found in the output directory")
		}
		if _, err := os.Stat(filepath.Join(outputPath, "b.txt")); !os.IsNotExist(err) {
			t.Error("unexpected file 'b.txt' was found in the output directory")
		}
	})

	t.Run("Invalid runs return errors", func(t *testing.T) {
		testCases := []struct {
			name string
			args []string
		}{
			{
				name: "Source does not exist",
				args: []string{flagSource, filepath.Join(rootDir, "non-existent"), flagManifest, filepath.Join(rootDir, testFileManifest), flagOutput, outputPath},
			},
			{
				name: "Source is a file, not a directory",
				args: []string{flagSource, filepath.Join(rootDir, testFileSource), flagManifest, filepath.Join(rootDir, testFileManifest), flagOutput, outputPath},
			},
			{
				name: "Manifest does not exist",
				args: []string{flagSource, filepath.Join(rootDir, testDirSource), flagManifest, filepath.Join(rootDir, "non-existent.txt"), flagOutput, outputPath},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := run(tc.args)
				if err == nil {
					t.Errorf("run() with args %v; expected an error but got nil", tc.args)
				}
			})
		}
	})
}
