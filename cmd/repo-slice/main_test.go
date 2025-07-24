// file: cmd/repo-slice/main_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestFS creates a temporary directory structure for testing.
// It returns the root directory path and a cleanup function.
func setupTestFS(t *testing.T) (string, func()) {
	t.Helper()
	rootDir, err := os.MkdirTemp("", "repo-slice-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create a valid source directory
	if err := os.Mkdir(filepath.Join(rootDir, "valid-source"), 0755); err != nil {
		t.Fatalf("failed to create source  dir: %v", err)
	}

	// Create a valid manifest file
	if err := os.WriteFile(filepath.Join(rootDir, "valid-manifest.txt"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create manifest file: %v", err)
	}

	// Create a file to act as an invalid source path
	if err := os.WriteFile(filepath.Join(rootDir, "source-is-a-file"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(rootDir)
	}

	return rootDir, cleanup
}

func TestRun_Validation(t *testing.T) {
	rootDir, cleanup := setupTestFS(t)
	defer cleanup()

	testCases := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Valid paths",
			args:    []string{"--source", filepath.Join(rootDir, "valid-source"), "--manifest", filepath.Join(rootDir, "valid-manifest.txt"), "--output", "test-output"},
			wantErr: false,
		},
		{
			name:    "Source does not exist",
			args:    []string{"--source", filepath.Join(rootDir, "non-existent"), "--manifest", filepath.Join(rootDir, "valid-manifest.txt"), "--output", "test-output"},
			wantErr: true,
		},
		{
			name:    "Source is a file, not a directory",
			args:    []string{"--source", filepath.Join(rootDir, "source-is-a-file"), "--manifest", filepath.Join(rootDir, "valid-manifest.txt"), "--output", "test-output"},
			wantErr: true,
		},
		{
			name:    "Manifest does not exist",
			args:    []string{"--source", filepath.Join(rootDir, "valid-source"), "--manifest", filepath.Join(rootDir, "non-existent.txt"), "--output", "test-output"},
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