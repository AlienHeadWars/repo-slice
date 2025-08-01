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

func TestSliceWithBasicIncludeAndExclude(t *testing.T) {
	sourceDir, outputDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Setup source files
	_ = os.WriteFile(filepath.Join(sourceDir, "fileA.txt"), []byte("a"), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "fileB.txt"), []byte("b"), 0644)
	_ = os.Mkdir(filepath.Join(sourceDir, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "docs", "guide.md"), []byte("guide"), 0644)

	// Setup manifest
	manifestPath := filepath.Join(sourceDir, "manifest.txt")
	manifestContent := "+\t/fileA.txt\n+\t/docs/\n-\t/docs/guide.md\n-\t*"
	_ = os.WriteFile(manifestPath, []byte(manifestContent), 0644)

	// Execute Slice
	err := Slice(sourceDir, outputDir, manifestPath, &CmdExecutor{})
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assertions
	assertFileExists(t, filepath.Join(outputDir, "fileA.txt"))
	assertFileExists(t, filepath.Join(outputDir, "docs"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "fileB.txt"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "docs", "guide.md"))
}

func TestSliceWithManifestInheritance(t *testing.T) {
	sourceDir, outputDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Setup source files
	_ = os.WriteFile(filepath.Join(sourceDir, "common.txt"), []byte("common"), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "feature.txt"), []byte("feature"), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test"), 0644)

	// Setup manifests
	baseManifestPath := filepath.Join(sourceDir, "base.manifest")
	_ = os.WriteFile(baseManifestPath, []byte("+\t/common.txt"), 0644)

	featureManifestPath := filepath.Join(sourceDir, "feature.manifest")
	_ = os.WriteFile(featureManifestPath, []byte(". base.manifest\n+\t/feature.txt\n-\t*"), 0644)

	// Execute Slice
	err := Slice(sourceDir, outputDir, featureManifestPath, &CmdExecutor{})
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assertions
	assertFileExists(t, filepath.Join(outputDir, "common.txt"))
	assertFileExists(t, filepath.Join(outputDir, "feature.txt"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "test.txt"))
}

func TestSliceWithWildcardInclusion(t *testing.T) {
	sourceDir, outputDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Setup source files
	_ = os.WriteFile(filepath.Join(sourceDir, "README.md"), []byte(""), 0644)
	_ = os.MkdirAll(filepath.Join(sourceDir, "pkg", "api"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "pkg", "api", "README.md"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "main.go"), []byte(""), 0644)

	// Setup manifest
	manifestPath := filepath.Join(sourceDir, "manifest.txt")
	manifestContent := "+\t**/README.md\n-\t*"
	_ = os.WriteFile(manifestPath, []byte(manifestContent), 0644)

	// Execute Slice
	err := Slice(sourceDir, outputDir, manifestPath, &CmdExecutor{})
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assertions
	assertFileExists(t, filepath.Join(outputDir, "README.md"))
	assertFileExists(t, filepath.Join(outputDir, "pkg", "api", "README.md"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "main.go"))
}

func TestSliceWithWildcardExclusion(t *testing.T) {
	sourceDir, outputDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Setup source files
	_ = os.WriteFile(filepath.Join(sourceDir, "app.go"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "app_test.go"), []byte(""), 0644)
	_ = os.Mkdir(filepath.Join(sourceDir, "logs"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "logs", "trace.log"), []byte(""), 0644)

	// Setup manifest
	manifestPath := filepath.Join(sourceDir, "manifest.txt")
	manifestContent := "+\t**\n-\t*.log\n-\t*_test.go"
	_ = os.WriteFile(manifestPath, []byte(manifestContent), 0644)

	// Execute Slice
	err := Slice(sourceDir, outputDir, manifestPath, &CmdExecutor{})
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assertions
	assertFileExists(t, filepath.Join(outputDir, "app.go"))
	assertFileExists(t, filepath.Join(outputDir, "logs"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "app_test.go"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "logs", "trace.log"))
}

func TestSliceWithRulePrecedence(t *testing.T) {
	sourceDir, outputDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Setup source files
	_ = os.MkdirAll(filepath.Join(sourceDir, "src"), 0755)
	_ = os.WriteFile(filepath.Join(sourceDir, "src", "core.go"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(sourceDir, "src", "utils.go"), []byte(""), 0644)

	// Setup manifest
	manifestPath := filepath.Join(sourceDir, "manifest.txt")
	manifestContent := "+\t/src/utils.go\n-\t/src/*"
	_ = os.WriteFile(manifestPath, []byte(manifestContent), 0644)

	// Execute Slice
	err := Slice(sourceDir, outputDir, manifestPath, &CmdExecutor{})
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assertions
	assertFileExists(t, filepath.Join(outputDir, "src", "utils.go"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "src", "core.go"))
}

func TestSliceWithSelfExclusion(t *testing.T) {
	sourceDir, outputDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Setup source files
	_ = os.WriteFile(filepath.Join(sourceDir, "file.txt"), []byte(""), 0644)
	manifestPath := filepath.Join(sourceDir, "the.manifest")
	manifestContent := "+\t/file.txt\n-\t/the.manifest"
	_ = os.WriteFile(manifestPath, []byte(manifestContent), 0644)

	// Execute Slice
	err := Slice(sourceDir, outputDir, manifestPath, &CmdExecutor{})
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assertions
	assertFileExists(t, filepath.Join(outputDir, "file.txt"))
	assertFileDoesNotExist(t, filepath.Join(outputDir, "the.manifest"))
}
