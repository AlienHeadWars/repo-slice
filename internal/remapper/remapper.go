// file: internal/remapper/remapper.go

// Package remapper provides functionality for renaming files in a directory
// based on a given extension map.
package remapper

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem defines an interface for file system operations needed by the
// remapper. This allows for a mock implementation in unit tests.
type FileSystem interface {
	WalkDir(root string, fn fs.WalkDirFunc) error
	Rename(oldpath, newpath string) error
}

// LiveFS is a concrete implementation of the FileSystem interface that uses
// the standard library's os and filepath packages.
type LiveFS struct{}

// WalkDir walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
func (fs *LiveFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

// Rename renames (moves) oldpath to newpath.
func (fs *LiveFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// ParseExtensionMap parses a comma-separated string of old:new pairs into a
// map of extensions to be remapped.
func ParseExtensionMap(mapStr string) (map[string]string, error) {
	if mapStr == "" {
		return map[string]string{}, nil
	}

	extMap := make(map[string]string)
	pairs := strings.Split(mapStr, ",")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed extension map pair: %s", pair)
		}

		oldExt := strings.TrimSpace(parts[0])
		newExt := strings.TrimSpace(parts[1])

		// Ensure extensions start with a dot for consistency.
		if !strings.HasPrefix(oldExt, ".") {
			oldExt = "." + oldExt
		}
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}

		extMap[oldExt] = newExt
	}

	return extMap, nil
}

// RemapExtensions walks a directory and renames files based on the provided
// extension map.
func RemapExtensions(dir string, extMap map[string]string, fsys FileSystem) error {
	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // Skip directories.
		}

		currentExt := filepath.Ext(path)
		newExt, shouldRemap := extMap[currentExt]

		if shouldRemap {
			base := strings.TrimSuffix(path, currentExt)
			newPath := base + newExt
			if err := fsys.Rename(path, newPath); err != nil {
				return fmt.Errorf("failed to rename %s to %s: %w", path, newPath, err)
			}
		}
		return nil
	}

	return fsys.WalkDir(dir, walkFn)
}
