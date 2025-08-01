// file: internal/slicer/slicer_test.go
package slicer

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// mockExecutor implements the Executor interface to capture command calls.
type mockExecutor struct {
	returnErr bool
	workDir   string
	command   string
	args      []string
}

func (m *mockExecutor) Run(workDir, command string, args ...string) error {
	m.workDir = workDir
	m.command = command
	m.args = args
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

// TestSliceConstructsCorrectFilterArgument verifies that the Slice function calls
// the executor with the correct rsync filter argument.
func TestSliceConstructsCorrectFilterArgument(t *testing.T) {
	sourceDir := "/source"
	outputDir := "/output"
	manifestPath := "/manifest.txt"
	mockExec := &mockExecutor{}

	err := Slice(sourceDir, outputDir, manifestPath, mockExec)
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	expectedFilterArg := fmt.Sprintf("merge %s", manifestPath)
	found := false
	for _, arg := range mockExec.args {
		if arg == expectedFilterArg {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("rsync arguments did not contain the correct filter rule.\nGot: %v\nWant rule: %s", mockExec.args, expectedFilterArg)
	}
}

func TestSlice(t *testing.T) {
	sourceDir := "/path/to/source"
	outputDir := "/tmp/output"
	manifestPath := "/path/to/manifest.txt"

	testCases := []struct {
		name    string
		exec    Executor
		wantErr bool
	}{
		{
			name:    "Happy path",
			exec:    &mockExecutor{},
			wantErr: false,
		},
		{
			name:    "Executor fails",
			exec:    &mockExecutor{returnErr: true},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Slice(sourceDir, outputDir, manifestPath, tc.exec)
			if (err != nil) != tc.wantErr {
				t.Errorf("Slice() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}