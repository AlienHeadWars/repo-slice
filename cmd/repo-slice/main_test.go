// file: cmd/repo-slice/main_test.go
package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

const (
	flagSource   = "--source"
	flagManifest = "--manifest"
	flagOutput   = "--output"
)

// mockFS is a mock implementation of the FileSystem interface for testing.
type mockFS struct {
	validateErr error
}

func (m *mockFS) ValidateInputs(cfg validate.Config) error { return m.validateErr }

// mockSlicer is a mock implementation of the Slicer interface for testing.
type mockSlicer struct {
	sliceErr error
}

func (m *mockSlicer) Slice(source, output, manifestPath string) error { return m.sliceErr }

// mockRemapper is a mock implementation of the Remapper interface for testing.
type mockRemapper struct {
	parseErr error
	remapErr error
}

func (m *mockRemapper) ParseExtensionMap(mapStr string) (map[string]string, error) {
	return nil, m.parseErr
}
func (m *mockRemapper) RemapExtensions(dir string, extMap map[string]string) error {
	return m.remapErr
}

// TestRunUnit tests the error-handling paths of the run function using mocks.
func TestRunUnit(t *testing.T) {
	validArgs := []string{flagManifest, "m.txt", flagSource, "s", flagOutput, "o"}
	remapArgs := []string{flagManifest, "m.txt", flagSource, "s", flagOutput, "o", "--extension-map", "tsx:ts"}

	testCases := []struct {
		name     string
		args     []string
		fs       FileSystem
		slicer   Slicer
		remapper Remapper
		wantErr  bool
	}{
		{"Argument parsing fails", []string{"--bad-flag"}, &mockFS{}, &mockSlicer{}, &mockRemapper{}, true},
		{"Validation fails", validArgs, &mockFS{validateErr: errors.New("validation failed")}, &mockSlicer{}, &mockRemapper{}, true},
		{"Slice operation fails", validArgs, &mockFS{}, &mockSlicer{sliceErr: errors.New("slice failed")}, &mockRemapper{}, true},
		{"Remap parsing fails", remapArgs, &mockFS{}, &mockSlicer{}, &mockRemapper{parseErr: errors.New("remap parse failed")}, true},
		{"Remap operation fails", remapArgs, &mockFS{}, &mockSlicer{}, &mockRemapper{remapErr: errors.New("remap op failed")}, true},
		{"Successful run", validArgs, &mockFS{}, &mockSlicer{}, &mockRemapper{}, false},
		{"Successful run with remap", remapArgs, &mockFS{}, &mockSlicer{}, &mockRemapper{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args, tc.fs, tc.slicer, tc.remapper)
			if (err != nil) != tc.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestRunIntegration is a simple end-to-end test.
func TestRunIntegration(t *testing.T) {
	rootDir, err := os.MkdirTemp("", "repo-slice-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(rootDir); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	})

	sourceDir := filepath.Join(rootDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "component.tsx"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create component.tsx: %v", err)
	}

	manifestPath := filepath.Join(rootDir, "manifest.txt")
	// Using rsync filter syntax for the integration test
	if err := os.WriteFile(manifestPath, []byte("+ component.tsx\n- *"), 0644); err != nil {
		t.Fatalf("failed to create manifest file: %v", err)
	}

	outputPath := filepath.Join(rootDir, "output")

	args := []string{flagManifest, manifestPath, flagSource, sourceDir, flagOutput, outputPath, "--extension-map", "tsx:ts"}
	err = run(args, &liveFS{}, &liveSlicer{}, &liveRemapper{})

	if err != nil {
		t.Fatalf("run() failed on integration test: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputPath, "component.ts")); os.IsNotExist(err) {
		t.Error("expected file 'component.ts' was not found in the output directory")
	}
}
