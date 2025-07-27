//go:build integration

// This file contains integration tests that interact with the real file system.
// To run these tests, use the build tag 'integration':
// go test -v ./... -tags=integration

package validate

import (
	"os"
	"testing"
)

func TestLiveFS_Stat(t *testing.T) {
	// Create a temporary file to test against.
	tmpFile, err := os.CreateTemp("", "test-stat-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test.

	fsys := LiveFS{}

	// Test case 1: A file that exists.
	t.Run("file exists", func(t *testing.T) {
		info, err := fsys.Stat(tmpFile.Name())
		if err != nil {
			t.Errorf("Stat() with existing file failed: %v", err)
		}
		if info == nil {
			t.Error("Stat() returned nil FileInfo for existing file")
		}
	})

	// Test case 2: A file that does not exist.
	t.Run("file does not exist", func(t *testing.T) {
		_, err := fsys.Stat("non-existent-file-path")
		if err == nil {
			t.Error("Stat() with non-existent file did not return an error")
		}
	})
}
