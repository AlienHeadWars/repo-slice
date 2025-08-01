//go:build integration

// file: internal/slicer/slicer_filter_integration_test.go
package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	readmeMdFile   = "README.md"
	commonTxtFile  = "common.txt"
	mainGoFile     = "main.go"
	mainTestGoFile = "main_test.go"
	appGoFile      = "src/app/app.go"
	appTestGoFile  = "src/app/app_test.go"
	guideMdFile    = "docs/guide.md"
	traceLogFile   = "docs/trace.log"
)

// setupIntegrationTest creates a temporary source and output directory for testing.
// It returns the paths to these directories and a cleanup function.
func setupIntegrationTest(t *testing.T) (sourceDir, outputDir string, cleanup func()) {
	t.Helper()
	rootDir, err := os.MkdirTemp("", "slicer-filter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp root dir: %v", err)
	}

	sourceDir = filepath.Join(rootDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	outputDir = filepath.Join(rootDir, "output")
	if err := os.Mkdir(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	cleanup = func() {
		os.RemoveAll(rootDir)
	}

	return sourceDir, outputDir, cleanup
}

// assertFileExists checks if a file exists and fails the test if it doesn't.
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist, but it doesn't: %s", path)
	}
}

// assertFileDoesNotExist checks if a file does not exist and fails the test if it does.
func assertFileDoesNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to not exist, but it does: %s", path)
	}
}

// setupCommonTestStructure creates a temporary source directory with a standardized
// set of files and directories to be used across all filter integration tests.
// It returns the source and output directory paths, and a cleanup function.
func setupCommonTestStructure(t *testing.T) (sourceDir, outputDir string, cleanup func()) {
	t.Helper()
	sourceDir, outputDir, cleanup = setupIntegrationTest(t)

	// Create a diverse file structure to test against
	_ = os.WriteFile(filepath.Join(sourceDir, readmeMdFile), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, commonTxtFile), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, mainGoFile), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, mainTestGoFile), []byte{}, 0644)

	// Create nested directories
	_ = os.MkdirAll(filepath.Join(sourceDir, "src", "app"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, appGoFile), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, appTestGoFile), []byte{}, 0644)

	_ = os.MkdirAll(filepath.Join(sourceDir, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, guideMdFile), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, traceLogFile), []byte{}, 0644)

	return sourceDir, outputDir, cleanup
}

// cleanManifest takes a multi-line string, trims whitespace from each line,
// and filters out empty lines. This allows manifests in tests to be indented.
func cleanManifest(s string) string {
	lines := strings.Split(s, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}
	return strings.Join(cleanedLines, "\n")
}

func TestSliceWithFiltering(t *testing.T) {
	sourceDir, outputDir, cleanup := setupCommonTestStructure(t)
	defer cleanup()

	// Base manifest for inheritance tests
	baseManifestPath := filepath.Join(sourceDir, "base.manifest")
	_ = os.WriteFile(baseManifestPath, []byte("+ /common.txt"), 0644)

	testCases := []struct {
		name               string
		manifestContent    string
		manifestFilename   string
		expectedToExist    []string
		expectedToNotExist []string
	}{
		{
			name:             "Basic Include and Exclude",
			manifestFilename: "manifest1.txt",
			manifestContent: `
				+ /main.go
				+ /docs/
				- /docs/guide.md
				- *
			`,
			expectedToExist:    []string{mainGoFile, "docs"},
			expectedToNotExist: []string{commonTxtFile, guideMdFile},
		},
		{
			name:             "Manifest Inheritance",
			manifestFilename: "manifest2.txt",
			manifestContent: `
				. base.manifest
				+ /main.go
				- *
			`,
			expectedToExist:    []string{commonTxtFile, mainGoFile},
			expectedToNotExist: []string{readmeMdFile},
		},
		{
			name:             "Wildcard Inclusion",
			manifestFilename: "manifest3.txt",
			manifestContent: `
				+ **/*.md
				+ **/
				- *
			`,
			expectedToExist:    []string{readmeMdFile, guideMdFile},
			expectedToNotExist: []string{mainGoFile},
		},
		{
			name:             "Wildcard Exclusion",
			manifestFilename: "manifest4.txt",
			manifestContent: `
				- *.log
				- *_test.go
				+ **
			`,
			expectedToExist:    []string{mainGoFile, appGoFile},
			expectedToNotExist: []string{traceLogFile, mainTestGoFile, appTestGoFile},
		},
		{
			name:             "Rule Precedence",
			manifestFilename: "manifest5.txt",
			manifestContent: `
				- /src/app/app_test.go
				+ /src/**
			`,
			expectedToExist:    []string{appGoFile},
			expectedToNotExist: []string{appTestGoFile},
		},
		{
			name:             "Self Exclusion of Manifest",
			manifestFilename: "manifest6.txt",
			manifestContent: `
				+ /main.go
				- /manifest6.txt
			`,
			expectedToExist:    []string{mainGoFile},
			expectedToNotExist: []string{"manifest6.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifestPath := filepath.Join(sourceDir, tc.manifestFilename)
			// Use the helper to clean the manifest content before writing
			_ = os.WriteFile(manifestPath, []byte(cleanManifest(tc.manifestContent)), 0644)

			err := Slice(sourceDir, outputDir, manifestPath, &CmdExecutor{})
			if err != nil {
				t.Fatalf("Slice() returned an unexpected error: %v", err)
			}

			for _, file := range tc.expectedToExist {
				assertFileExists(t, filepath.Join(outputDir, file))
			}
			for _, file := range tc.expectedToNotExist {
				assertFileDoesNotExist(t, filepath.Join(outputDir, file))
			}

			// Clean the output directory for the next run
			os.RemoveAll(outputDir)
			os.Mkdir(outputDir, 0755)
		})
	}
}
