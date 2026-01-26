package main

import (
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	originalExecute := execute
	execute = func() error { return nil }
	t.Cleanup(func() { execute = originalExecute })

	if err := run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMainExitOnError(t *testing.T) {
	originalExecute := execute
	originalExit := exit
	executeErr := errors.New("execute failed")
	exitCode := 0

	execute = func() error { return executeErr }
	exit = func(code int) { exitCode = code }
	t.Cleanup(func() {
		execute = originalExecute
		exit = originalExit
	})

	main()
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}
