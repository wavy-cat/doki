package cmd

import (
	"os"
	"testing"
)

func TestExecuteHelp(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"doki", "--help"}
	t.Cleanup(func() { os.Args = originalArgs })

	if err := Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
