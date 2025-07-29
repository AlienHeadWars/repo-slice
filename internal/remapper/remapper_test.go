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
	baseExtMap := map[string]string{".tsx": ".ts"}

	testCases := []struct {
		name           string
		mockFS         *mocks.MockFS
		extMap         map[string]string
		wantErr        bool
		expectedRename string
	}{
		{
			name:           "renames matching file",
			mockFS:         &mocks.MockFS{Files: map[string]bool{"component.tsx": false}},
			extMap:         baseExtMap,
			wantErr:        false,
			expectedRename: "component.ts",
		},
		{
			name:           "ignores non-matching file",
			mockFS:         &mocks.MockFS{Files: map[string]bool{"style.css": false}},
			extMap:         baseExtMap,
			wantErr:        false,
			expectedRename: "",
		},
		{
			name: "handles rename error",
			mockFS: &mocks.MockFS{
				Files:     map[string]bool{"component.tsx": false},
				RenameErr: errors.New("rename failed"),
			},
			extMap:         baseExtMap,
			wantErr:        true,
			expectedRename: "component.ts", // It will attempt to rename to this
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RemapExtensions(".", tc.extMap, tc.mockFS)

			if (err != nil) != tc.wantErr {
				t.Fatalf("RemapExtensions() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.mockFS.RenameTo != tc.expectedRename {
				t.Errorf("Expected rename to '%s', got '%s'", tc.expectedRename, tc.mockFS.RenameTo)
			}
		})
	}
}
