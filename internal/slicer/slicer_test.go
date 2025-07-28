// file: internal/slicer/slicer_test.go
package slicer

import (
	"errors"
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
	returnErr bool // New field to make the mock return an error.
}

func (m *mockExecutor) Run(workDir, command string, args ...string) error {
	m.calledWith.workDir = workDir
	m.calledWith.command = command
	m.calledWith.args = args
	if m.returnErr {
		return errors.New("mock executor error")
	}
	return nil
}

// errorReader is a mock io.Reader that always returns an error.
type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock reader error")
}

func TestParseManifest(t *testing.T) {
	t.Run("successful parse", func(t *testing.T) {
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
	})

	t.Run("reader error", func(t *testing.T) {
		_, err := ParseManifest(errorReader{})
		if err == nil {
			t.Error("ParseManifest() with a failing reader should have returned an error, but did not")
		}
	})
}

func TestSlice(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
	})

	t.Run("executor error", func(t *testing.T) {
		sourceDir := "/path/to/source"
		outputDir := "/tmp/output"
		files := []string{"file1.go"}

		executor := &mockExecutor{returnErr: true}
		err := Slice(sourceDir, outputDir, files, executor)

		if err == nil {
			t.Error("Slice() with a failing executor should have returned an error, but did not")
		}
	})
}
