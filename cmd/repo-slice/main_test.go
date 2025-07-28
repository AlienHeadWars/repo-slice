// file: cmd/repo-slice/main_test.go
package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AlienHeadwars/repo-slice/internal/validate"
)

// mockFS is a mock implementation of the FileSystem interface for testing.
type mockFS struct {
	validateErr error
	openErr     error
}

func (m *mockFS) ValidateInputs(cfg validate.Config) error { return m.validateErr }
func (m *mockFS) Open(name string) (io.ReadCloser, error) {
	if m.openErr != nil {
		return nil, m.openErr
	}
	// Return a no-op closer for the success path.
	return io.NopCloser(strings.NewReader("")), nil
}

// mockSlicer is a mock implementation of the Slicer interface for testing.
type mockSlicer struct {
	parseErr error
	sliceErr error
}

func (m *mockSlicer) ParseManifest(r io.Reader) ([]string, error)       { return nil, m.parseErr }
func (m *mockSlicer) Slice(source, output string, files []string) error { return m.sliceErr }

// TestRunUnit tests the error-handling paths of the run function using mocks.
func TestRunUnit(t *testing.T) {
	// Dummy args for tests that get past the parsing stage.
	validArgs := []string{"--manifest", "m.txt", "--source", "s", "--output", "o"}

	testCases := []struct {
		name    string
		args    []string
		fs      FileSystem
		slicer  Slicer
		wantErr bool
	}{
		{"Argument parsing fails", []string{"--bad-flag"}, &mockFS{}, &mockSlicer{}, true},
		{"Validation fails", validArgs, &mockFS{validateErr: errors.New("validation failed")}, &mockSlicer{}, true},
		{"File open fails", validArgs, &mockFS{openErr: errors.New("open failed")}, &mockSlicer{}, true},
		{"Manifest parsing fails", validArgs, &mockFS{}, &mockSlicer{parseErr: errors.New("parse failed")}, true},
		{"Slice operation fails", validArgs, &mockFS{}, &mockSlicer{sliceErr: errors.New("slice failed")}, true},
		{"Successful run", validArgs, &mockFS{}, &mockSlicer{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args, tc.fs, tc.slicer)
			if (err != nil) != tc.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestRunIntegration is a simple end-to-end test to ensure the real
// components are wired together correctly for the happy path.
func TestRunIntegration(t *testing.T) {
	rootDir, err := os.MkdirTemp("", "repo-slice-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir)

	sourceDir := filepath.Join(rootDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "a.txt"), []byte("a"), 0644)

	manifestPath := filepath.Join(rootDir, "manifest.txt")
	os.WriteFile(manifestPath, []byte("a.txt"), 0644)

	outputPath := filepath.Join(rootDir, "output")

	args := []string{"--manifest", manifestPath, "--source", sourceDir, "--output", outputPath}
	err = run(args, &liveFS{}, &liveSlicer{})

	if err != nil {
		t.Fatalf("run() failed on integration test: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputPath, "a.txt")); os.IsNotExist(err) {
		t.Error("expected file 'a.txt' was not found in the output directory")
	}
}
