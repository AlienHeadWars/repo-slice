// file: internal/validate/validate_test.go
package validate

import (
	"testing"

	"github.com/AlienHeadwars/repo-slice/internal/mocks"
)

const (
	validSource     = "valid-source"
	validManifest   = "valid-manifest.txt"
	sourceAsFile    = "path-is-a-file"
	manifestAsDir   = "path-is-a-dir"
	nonExistentFile = "non-existent"
	nonExistentTxt  = "non-existent.txt"
)

// TestValidateInputs is a unit test that checks the input validation logic.
func TestValidateInputs(t *testing.T) {
	fsys := &mocks.MockFS{
		Files: map[string]bool{
			validSource:   true,
			validManifest: false,
			sourceAsFile:  false,
			manifestAsDir: true,
		},
	}

	testCases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"Valid paths", Config{SourcePath: validSource, ManifestPath: validManifest}, false},
		{"Source does not exist", Config{SourcePath: nonExistentFile, ManifestPath: validManifest}, true},
		{"Source is a file", Config{SourcePath: sourceAsFile, ManifestPath: validManifest}, true},
		{"Manifest does not exist", Config{SourcePath: validSource, ManifestPath: nonExistentTxt}, true},
		{"Manifest is a directory", Config{SourcePath: validSource, ManifestPath: manifestAsDir}, true},
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
