// file: internal/slicer/slicer_test.go
package slicer

import (
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"time"
)

// ... (mockFileInfo and mockExecutor structs remain the same) ...
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

type mockExecutor struct {
	calledWith struct {
		command string
		args    []string
	}
}

func (m *mockExecutor) Run(command string, args ...string) error {
	m.calledWith.command = command
	m.calledWith.args = args
	return nil
}

func TestParseManifest(t *testing.T) {
	manifestContent := `
		file1.go
		  dir/file2.go  
		
		# A comment to be ignored by the logic
		dir2/
	`
	// Use strings.NewReader to create an io.Reader from the test string.
	reader := strings.NewReader(manifestContent)

	expected := []string{"file1.go", "dir/file2.go", "# A comment to be ignored by the logic", "dir2/"}
	actual, err := ParseManifest(reader)

	if err != nil {
		t.Fatalf("ParseManifest() returned an unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ParseManifest() = %v, want %v", actual, expected)
	}
}

func TestSlice(t *testing.T) {
	sourceDir := "/path/to/source"
	outputDir := "/tmp/output"
	files := []string{"file1.go", "dir/file2.go"}

	executor := &mockExecutor{}
	err := Slice(sourceDir, outputDir, files, executor)

	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	if executor.calledWith.command != "rsync" {
		t.Errorf("Expected command to be 'rsync', but got '%s'", executor.calledWith.command)
	}

	argsString := strings.Join(executor.calledWith.args, " ")
	if !strings.Contains(argsString, "-a") {
		t.Error("rsync arguments should contain '-a'")
	}
	if !strings.Contains(argsString, "--files-from=") {
		t.Error("rsync arguments should contain '--files-from='")
	}
	if !strings.Contains(argsString, sourceDir) {
		t.Errorf("rsync arguments should contain source directory '%s'", sourceDir)
	}
	if !strings.Contains(argsString, outputDir) {
		t.Errorf("rsync arguments should contain output directory '%s'", outputDir)
	}
}
