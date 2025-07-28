// file: internal/slicer/slicer.go

// Package slicer contains the core logic for parsing the manifest and
// executing the rsync command to create a repository slice.
package slicer

// Executor defines an interface for running external commands. This allows for
// a mock implementation to be used in unit tests, adhering to the Dependency
// Inversion Principle.
type Executor interface {
	// Run executes a command with the given arguments and returns an error
	// if the command fails.
	Run(command string, args ...string) error
}