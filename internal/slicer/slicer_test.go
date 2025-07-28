// file: internal/slicer/slicer_test.go
package slicer

import (
	"reflect"
	"strings"
	"testing"
)

// mockExecutor implements the Executor interface to capture command calls.
type mockExecutor struct {
	calledWith struct {
		workDir string
		command string
		args    []string
	}
}

func (m *mockExecutor) Run(workDir, command string, args ...string) error {
	m.calledWith.workDir = workDir
	m.calledWith.command = command
	m.calledWith.args = args
	return nil
}

func TestParseManifest(t *testing.T) {
	manifestContent := `
		file1.go
		  dir/file2.go  
		# A comment to be ignored
		dir2/
	`
	reader := strings.NewReader(manifestContent)

	expected := []string{"file1.go", "dir/file2.go", "dir2/"}
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

	if executor.calledWith.workDir != sourceDir {
		t.Errorf("Expected workDir to be '%s', got '%s'", sourceDir, executor.calledWith.workDir)
	}
	if executor.calledWith.command != "rsync" {
		t.Errorf("Expected command to be 'rsync', but got '%s'", executor.calledWith.command)
	}
}
