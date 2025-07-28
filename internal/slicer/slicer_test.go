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
	returnErr bool
}

func (m *mockExecutor) Run(workDir, command string, args ...string) error {
	if m.returnErr {
		return errors.New("mock executor error")
	}
	return nil
}

// mockTempFile implements the TempFile interface for testing.
type mockTempFile struct {
	writeErr error
	closeErr error
}

func (m *mockTempFile) Write(p []byte) (n int, err error) { return 0, m.writeErr }
func (m *mockTempFile) Close() error                      { return m.closeErr }
func (m *mockTempFile) Name() string                      { return "tempfile" }

// mockTempFiler implements the TempFiler interface for testing.
type mockTempFiler struct {
	createErr error
	file      *mockTempFile
}

func (m *mockTempFiler) CreateTemp(dir, pattern string) (TempFile, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.file, nil
}

// errorReader is a mock io.Reader that always returns an error.
type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock reader error")
}

func TestParseManifest(t *testing.T) {
	// ... (This test remains the same)
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
	sourceDir := "/path/to/source"
	outputDir := "/tmp/output"
	files := []string{"file1.go"}

	testCases := []struct {
		name    string
		filer   TempFiler
		exec    Executor
		wantErr bool
	}{
		{
			name:    "Happy path",
			filer:   &mockTempFiler{file: &mockTempFile{}},
			exec:    &mockExecutor{},
			wantErr: false,
		},
		{
			name:    "Temp file creation fails",
			filer:   &mockTempFiler{createErr: errors.New("create failed")},
			exec:    &mockExecutor{},
			wantErr: true,
		},
		{
			name:    "Temp file write fails",
			filer:   &mockTempFiler{file: &mockTempFile{writeErr: errors.New("write failed")}},
			exec:    &mockExecutor{},
			wantErr: true,
		},
		{
			name:    "Temp file close fails",
			filer:   &mockTempFiler{file: &mockTempFile{closeErr: errors.New("close failed")}},
			exec:    &mockExecutor{},
			wantErr: true,
		},
		{
			name:    "Executor fails",
			filer:   &mockTempFiler{file: &mockTempFile{}},
			exec:    &mockExecutor{returnErr: true},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Slice(sourceDir, outputDir, files, tc.exec, tc.filer)
			if (err != nil) != tc.wantErr {
				t.Errorf("Slice() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
