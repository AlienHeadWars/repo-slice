//go:build integration

// This file contains integration tests for the remapper package that interact
// with the real file system. To run these tests, use the build tag 'integration':
// go test -v ./... -tags=integration

package remapper

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestLiveFSWalkAndRename(t *testing.T) {
	// Setup a temporary directory with a file.
	rootDir, err := os.MkdirTemp("", "livefs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir)

	oldPath := filepath.Join(rootDir, "oldname.txt")
	newPath := filepath.Join(rootDir, "newname.txt")

	if err := os.WriteFile(oldPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	fsys := &LiveFS{}

	// Test WalkDir by verifying it finds our file.
	found := false
	walkFn := func(path string, d fs.DirEntry, err error) error {
		if path == oldPath {
			found = true
		}

		return nil
	}
	if err := fsys.WalkDir(rootDir, walkFn); err != nil {
		t.Fatalf("WalkDir failed: %v", err)
	}
	if !found {
		t.Error("WalkDir did not find the test file")
	}

	// Test Rename.
	if err := fsys.Rename(oldPath, newPath); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	// Verify the rename was successful.
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Error("Rename failed: new file path does not exist")
	}
}
