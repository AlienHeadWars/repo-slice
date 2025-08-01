//go:build integration

// file: internal/slicer/slicer_filter_integration_test.go
package slicer

import (
	"os"
	"path/filepath"
	"testing"
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
	_ = os.WriteFile(filepath.Join(sourceDir, "README.md"), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "common.txt"), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "main.go"), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "main_test.go"), []byte{}, 0644)

	// Create nested directories
	_ = os.MkdirAll(filepath.Join(sourceDir, "src", "app"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "src", "app", "app.go"), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "src", "app", "app_test.go"), []byte{}, 0644)

	_ = os.MkdirAll(filepath.Join(sourceDir, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "docs", "guide.md"), []byte{}, 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "docs", "trace.log"), []byte{}, 0644)

	return sourceDir, outputDir, cleanup
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
			manifestContent: `+ /main.go
+ /docs/
- /docs/guide.md
- *
`,
			// Expects main.go and the docs dir to be included, but guide.md
			// within docs is explicitly excluded. common.txt is excluded by the final '- *'.
			expectedToExist:    []string{"main.go", "docs"},
			expectedToNotExist: []string{"common.txt", "docs/guide.md"},
		},
		{
			name:             "Manifest Inheritance",
			manifestFilename: "manifest2.txt",
			manifestContent: `. base.manifest
+ /main.go
- *
`,
			// Expects common.txt (from base.manifest) and main.go to be included.
			// README.md is excluded by the final '- *'.
			expectedToExist:    []string{"common.txt", "main.go"},
			expectedToNotExist: []string{"README.md"},
		},
		{
			name:             "Wildcard Inclusion",
			manifestFilename: "manifest3.txt",
			manifestContent: `+ **/*.md
- *
`,
			// Expects both README.md and the nested docs/guide.md to be included
			// due to the recursive wildcard match.
			expectedToExist:    []string{"README.md", "docs/guide.md"},
			expectedToNotExist: []string{"main.go"},
		},
		{
			name:             "Wildcard Exclusion",
			manifestFilename: "manifest4.txt",
			manifestContent: `+ **
- *.log
- *_test.go
`,
			// Expects everything to be included except for files ending in .log or _test.go.
			expectedToExist:    []string{"main.go", "src/app/app.go"},
			expectedToNotExist: []string{"docs/trace.log", "main_test.go", "src/app/app_test.go"},
		},
		{
			name:             "Rule Precedence",
			manifestFilename: "manifest5.txt",
			manifestContent: `+ /src/app/app.go
- /src/**
`,
			// The more specific include rule for app.go should take precedence over the
			// general exclude rule for the src directory's contents.
			expectedToExist:    []string{"src/app/app.go"},
			expectedToNotExist: []string{},
		},
		{
			name:             "Self Exclusion of Manifest",
			manifestFilename: "manifest6.txt",
			manifestContent: `+ /main.go
- /manifest6.txt
`,
			// The manifest should include main.go but exclude itself from the output.
			expectedToExist:    []string{"main.go"},
			expectedToNotExist: []string{"manifest6.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifestPath := filepath.Join(sourceDir, tc.manifestFilename)
			_ = os.WriteFile(manifestPath, []byte(tc.manifestContent), 0644)

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
