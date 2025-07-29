// file: internal/validate/validate_test.go
package validate

import (
	"testing"

	"github.com/AlienHeadwars/repo-slice/internal/mocks"
)

// TestValidateInputs is a unit test that checks the input validation logic.
func TestValidateInputs(t *testing.T) {
	fsys := &mocks.MockFS{
		Files: map[string]bool{
			"valid-source":       true,
			"valid-manifest.txt": false,
			"path-is-a-file":     false,
			"path-is-a-dir":      true,
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
