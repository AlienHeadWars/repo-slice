// file: internal/remapper/remapper_test.go
package remapper

import (
	"errors"
	"reflect"
	"testing"

	"github.com/AlienHeadwars/repo-slice/internal/mocks"
)

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
		fsys := &mocks.MockFS{Files: map[string]bool{"component.tsx": false}}
		extMap := map[string]string{".tsx": ".ts"}
		err := RemapExtensions(".", extMap, fsys)
		if err != nil {
			t.Fatalf("RemapExtensions() failed: %v", err)
		}
		if fsys.RenameTo != "component.ts" {
			t.Errorf("Expected rename to 'component.ts', got '%s'", fsys.RenameTo)
		}
	})

	t.Run("ignores non-matching file", func(t *testing.T) {
		fsys := &mocks.MockFS{Files: map[string]bool{"style.css": false}}
		extMap := map[string]string{".tsx": ".ts"}
		err := RemapExtensions(".", extMap, fsys)
		if err != nil {
			t.Fatalf("RemapExtensions() failed: %v", err)
		}
		if fsys.RenameTo != "" {
			t.Errorf("Expected no rename, but got rename to '%s'", fsys.RenameTo)
		}
	})

	t.Run("handles rename error", func(t *testing.T) {
		fsys := &mocks.MockFS{
			Files:     map[string]bool{"component.tsx": false},
			RenameErr: errors.New("rename failed"),
		}
		extMap := map[string]string{".tsx": ".ts"}
		err := RemapExtensions(".", extMap, fsys)
		if err == nil {
			t.Error("Expected an error from rename failure, but got nil")
		}
	})
}
