// file: internal/remapper/remapper_test.go
package remapper

import (
	"errors"
	"io/fs"
	"reflect"
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

// mockFS is a mock implementation of the FileSystem interface for testing.
type mockFS struct {
	files      map[string]bool // path -> isDir
	renameErr  error
	renameFrom string
	renameTo   string
}

func (m *mockFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	for path, isDir := range m.files {
		// A simplified mock of WalkDir for testing purposes.
		d := fs.FileInfoToDirEntry(mockFileInfo{name: path, isDir: isDir})
		if err := fn(path, d, nil); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockFS) Rename(oldpath, newpath string) error {
	m.renameFrom = oldpath
	m.renameTo = newpath
	return m.renameErr
}

func TestParseExtensionMap(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{"Valid map", "tsx:ts,mdx:md", map[string]string{".tsx": ".ts", ".mdx": ".md"}, false},
		{"Empty map", "", map[string]string{}, false},
		{"Whitespace handling", " tsx : ts , mdx:md ", map[string]string{".tsx": ".ts", ".mdx": ".md"}, false},
		{"Malformed pair", "tsx:ts,mdx", nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseExtensionMap(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseExtensionMap() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseExtensionMap() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRemapExtensions(t *testing.T) {
	t.Run("renames matching file", func(t *testing.T) {
		fsys := &mockFS{files: map[string]bool{"component.tsx": false}}
		extMap := map[string]string{".tsx": ".ts"}
		err := RemapExtensions(".", extMap, fsys)
		if err != nil {
			t.Fatalf("RemapExtensions() failed: %v", err)
		}
		if fsys.renameTo != "component.ts" {
			t.Errorf("Expected rename to 'component.ts', got '%s'", fsys.renameTo)
		}
	})

	t.Run("ignores non-matching file", func(t *testing.T) {
		fsys := &mockFS{files: map[string]bool{"style.css": false}}
		extMap := map[string]string{".tsx": ".ts"}
		err := RemapExtensions(".", extMap, fsys)
		if err != nil {
			t.Fatalf("RemapExtensions() failed: %v", err)
		}
		if fsys.renameTo != "" {
			t.Errorf("Expected no rename, but got rename to '%s'", fsys.renameTo)
		}
	})

	t.Run("handles rename error", func(t *testing.T) {
		fsys := &mockFS{
			files:     map[string]bool{"component.tsx": false},
			renameErr: errors.New("rename failed"),
		}
		extMap := map[string]string{".tsx": ".ts"}
		err := RemapExtensions(".", extMap, fsys)
		if err == nil {
			t.Error("Expected an error from rename failure, but got nil")
		}
	})
}
