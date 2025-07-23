// file: cmd/repo-slice/main_test.go
package main

import "testing"

func TestRun_ArgumentParsing(t *testing.T) {
	// Define the test arguments.
	args := []string{
		"--manifest", "test-manifest.txt",
		"--output", "test-output",
		"--source", "test-source",
	}

	// Run the function with the test arguments.
	// A nil error indicates successful parsing.
	if err := run(args); err != nil {
		t.Errorf("run() with valid arguments failed: %v", err)
	}
}