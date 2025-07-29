// file: internal/remapper/remapper.go

// Package remapper provides functionality for renaming files in a directory
// based on a given extension map.
package remapper

import (
	"io/fs"
)

// FileSystem defines an interface for file system operations needed by the
// remapper. This allows for a mock implementation in unit tests.
type FileSystem interface {
	WalkDir(root string, fn fs.WalkDirFunc) error
	Rename(oldpath, newpath string) error
}

// ParseExtensionMap parses a comma-separated string of old:new pairs into a
// map of extensions to be remapped.
func ParseExtensionMap(mapStr string) (map[string]string, error) {
	// TODO: Implement parsing logic.
	return nil, nil
}

// RemapExtensions walks a directory and renames files based on the provided
// extension map.
func RemapExtensions(dir string, extMap map[string]string, fsys FileSystem) error {
	// TODO: Implement remapping logic.
	return nil
}