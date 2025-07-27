// file: internal/validate/validate_test.go
package validate

import (
	"io/fs"
	"testing"
	"time"
)

const (
	validSourcePath      = "valid-source"
	validManifestPath    = "valid-manifest.txt"
	invalidSourceAsFile  = "path-is-a-file"
	invalidManifestAsDir = "path-is-a-dir"
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

// TestValidateInputs is a unit test that checks the input validation logic.
// It uses a mock file system to simulate various scenarios, such as missing
// files or incorrect path types, ensuring the validation function returns
// errors when expected.
func TestValidateInputs(t *testing.T) {
	fsys := mockFS{
		files: map[string]mockFileInfo{
			validSourcePath:      {name: validSourcePath, isDir: true},
			validManifestPath:    {name: validManifestPath, isDir: false},
			invalidSourceAsFile:  {name: invalidSourceAsFile, isDir: false},
			invalidManifestAsDir: {name: invalidManifestAsDir, isDir: true},
		},
	}

	testCases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"Valid paths", Config{SourcePath: validSourcePath, ManifestPath: validManifestPath}, false},
		{"Source does not exist", Config{SourcePath: "non-existent", ManifestPath: validManifestPath}, true},
		{"Source is a file", Config{SourcePath: invalidSourceAsFile, ManifestPath: validManifestPath}, true},
		{"Manifest does not exist", Config{SourcePath: validSourcePath, ManifestPath: "non-existent.txt"}, true},
		{"Manifest is a directory", Config{SourcePath: validSourcePath, ManifestPath: invalidManifestAsDir}, true},
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
