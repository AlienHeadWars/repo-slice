// file: internal/validate/validate_test.go
package validate

import (
	"io/fs"
	"testing"
	"time"
)

// mockFileInfo implements fs.FileInfo for our mock file system.
type mockFileInfo struct {
	name  string
	isDir bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

// mockFS implements the FS interface for testing purposes.
type mockFS struct {
	files map[string]mockFileInfo
}

// Stat simulates the Stat operation for our mock file system.
func (m mockFS) Stat(name string) (fs.FileInfo, error) {
	if file, ok := m.files[name]; ok {
		return file, nil
	}
	return nil, fs.ErrNotExist
}

func TestValidateInputs(t *testing.T) {
	// Setup a mock file system state
	fsys := mockFS{
		files: map[string]mockFileInfo{
			"valid-source":      {name: "valid-source", isDir: true},
			"valid-manifest.txt": {name: "valid-manifest.txt", isDir: false},
			"path-is-a-file":    {name: "path-is-a-file", isDir: false},
			"path-is-a-dir":     {name: "path-is-a-dir", isDir: true},
		},
	}

	testCases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"Valid paths", Config{SourcePath: "valid-source", ManifestPath: "valid-manifest.txt"}, false},
		{"Source does not exist", Config{SourcePath: "non-existent", ManifestPath: "valid-manifest.txt"}, true},
		{"Source is a file", Config{SourcePath: "path-is-a-file", ManifestPath: "valid-manifest.txt"}, true},
		{"Manifest does not exist", Config{SourcePath: "valid-source", ManifestPath: "non-existent.txt"}, true},
		{"Manifest is a directory", Config{SourcePath: "valid-source", ManifestPath: "path-is-a-dir"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateInputs(tc.cfg, fsys)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateInputs() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}