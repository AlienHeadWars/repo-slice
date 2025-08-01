//go:build integration

package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCmdExecutorRunFails is an integration test for the CmdExecutor that
// verifies its error handling for various failure scenarios.
func TestCmdExecutorRunFails(t *testing.T) {
	executor := CmdExecutor{}

	t.Run("invalid flag", func(t *testing.T) {
		err := executor.Run(".", "rsync", "--non-existent-flag")
		if err == nil {
			t.Fatal("CmdExecutor.Run() did not return an error for a failing command")
		}
		if !strings.Contains(err.Error(), "STDERR") {
			t.Error("Error message from failed command did not include STDERR")
		}
	})

	t.Run("permission denied", func(t *testing.T) {
		unwritableDir, err := os.MkdirTemp("", "unwritable-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(unwritableDir)
		if err := os.Chmod(unwritableDir, 0555); err != nil {
			t.Fatalf("Failed to chmod temp dir: %v", err)
		}

		err = executor.Run(".", "touch", filepath.Join(unwritableDir, "test.txt"))
		if err == nil {
			t.Fatal("CmdExecutor.Run() did not return an error for a permission denied scenario")
		}
	})

	t.Run("command succeeds but writes to stderr", func(t *testing.T) {
		// This test covers the specific case where a command exits with code 0,
		// but still produces output on stderr, which we treat as an error.
		errMsg := "this is a stderr message"
		err := executor.Run(".", "sh", "-c", "echo '"+errMsg+"' >&2; exit 0")

		if err == nil {
			t.Fatal("CmdExecutor.Run() did not return an error for a successful command that wrote to stderr")
		}
		if !strings.Contains(err.Error(), errMsg) {
			t.Errorf("Error message did not contain the expected stderr output. Got: %v", err)
		}
	})
}
