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

// liveFS is a concrete implementation of the FileSystem interface.
type liveFS struct{}

func (fs *liveFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

func (fs *liveFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func TestRemapExtensions_Integration(t *testing.T) {
	// Setup a temporary directory with files to be renamed.
	rootDir, err := os.MkdirTemp("", "remap-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir)

	_ = os.WriteFile(filepath.Join(rootDir, "component.tsx"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(rootDir, "style.css"), []byte(""), 0644)

	extMap := map[string]string{".tsx": ".ts"}
	fsys := &liveFS{}

	// Execute the remapping.
	err = RemapExtensions(rootDir, extMap, fsys)
	if err != nil {
		t.Fatalf("RemapExtensions() failed: %v", err)
	}

	// Assert that the files have been correctly renamed.
	if _, err := os.Stat(filepath.Join(rootDir, "component.ts")); os.IsNotExist(err) {
		t.Error("Expected 'component.tsx' to be renamed to 'component.ts', but it was not found.")
	}
	if _, err := os.Stat(filepath.Join(rootDir, "component.tsx")); !os.IsNotExist(err) {
		t.Error("Expected 'component.tsx' to be removed after rename, but it still exists.")
	}
	if _, err := os.Stat(filepath.Join(rootDir, "style.css")); os.IsNotExist(err) {
		t.Error("Expected 'style.css' to be untouched, but it was not found.")
	}
}
