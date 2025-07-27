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

	if err := os.Mkdir(filepath.Join(rootDir, testDirSource), 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(rootDir, testFileManifest), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create manifest file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(rootDir, testFileSource), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	return rootDir
}

// TestRunValidation is an integration test for the main run function. It verifies
// that the application correctly handles valid and invalid command-line
// arguments by checking if errors are returned appropriately.
func TestRunValidation(t *testing.T) {
	rootDir := setupTestFS(t)

	testCases := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Valid paths",
			args:    []string{flagSource, filepath.Join(rootDir, testDirSource), flagManifest, filepath.Join(rootDir, testFileManifest), flagOutput, testFileOutput},
			wantErr: false,
		},
		{
			name:    "Source does not exist",
			args:    []string{flagSource, filepath.Join(rootDir, "non-existent"), flagManifest, filepath.Join(rootDir, testFileManifest), flagOutput, testFileOutput},
			wantErr: true,
		},
		{
			name:    "Source is a file, not a directory",
			args:    []string{flagSource, filepath.Join(rootDir, testFileSource), flagManifest, filepath.Join(rootDir, testFileManifest), flagOutput, testFileOutput},
			wantErr: true,
		},
		{
			name:    "Manifest does not exist",
			args:    []string{flagSource, filepath.Join(rootDir, testDirSource), flagManifest, filepath.Join(rootDir, "non-existent.txt"), flagOutput, testFileOutput},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args)
			if (err != nil) != tc.wantErr {
				t.Errorf("run() with args %v; got error = %v, wantErr %v", tc.args, err, tc.wantErr)
			}
		})
	}
}
