// file: internal/mocks/mocks.go

// Package mocks provides shared mock implementations for testing across the application.
package mocks

import (
	"io/fs"
	"time"
)

// MockFileInfo implements fs.FileInfo for our mock file system.
type MockFileInfo struct {
	FileName  string
	IsDirBool bool
}

// Name returns the name of the file.
func (m MockFileInfo) Name() string { return m.FileName }

// Size returns the size of the file.
func (m MockFileInfo) Size() int64 { return 0 }

// Mode returns the file mode.
func (m MockFileInfo) Mode() fs.FileMode { return 0 }

// ModTime returns the modification time.
func (m MockFileInfo) ModTime() time.Time { return time.Time{} }

// IsDir returns true if the entry is a directory.
func (m MockFileInfo) IsDir() bool { return m.IsDirBool }

// Sys returns the underlying data source.
func (m MockFileInfo) Sys() interface{} { return nil }

// MockFS implements the FS interface for testing purposes.
type MockFS struct {
	Files      map[string]bool // path -> isDir
	RenameErr  error
	WalkErr    error
	WalkFnErr  error // New field to simulate an error passed to the walk function.
	RenameFrom string
	RenameTo   string
}

// Stat simulates the Stat operation for our mock file system.
func (m *MockFS) Stat(name string) (fs.FileInfo, error) {
	isDir, ok := m.Files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return MockFileInfo{FileName: name, IsDirBool: isDir}, nil
}

// WalkDir simulates walking a directory structure.
func (m *MockFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	if m.WalkErr != nil {
		return m.WalkErr
	}
	for path, isDir := range m.Files {
		d := fs.FileInfoToDirEntry(MockFileInfo{FileName: path, IsDirBool: isDir})
		// Pass the WalkFnErr to the callback to simulate a file system error during iteration.
		if err := fn(path, d, m.WalkFnErr); err != nil {
			return err
		}
	}
	return nil
}

// Rename simulates renaming a file.
func (m *MockFS) Rename(oldpath, newpath string) error {
	m.RenameFrom = oldpath
	m.RenameTo = newpath
	return m.RenameErr
}
