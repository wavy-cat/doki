package cmd

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/wavy-cat/doki/internal/app"
)

func TestRootCommandExecute(t *testing.T) {
	t.Run("parses flags and runs", func(t *testing.T) {
		var (
			calls       int
			lastHost    string
			lastPort    uint16
			lastTimeout time.Duration
		)
		runner := app.Runner{
			ScanPort: func(host string, port uint16, timeout time.Duration) error {
				calls++
				lastHost = host
				lastPort = port
				lastTimeout = timeout
				return nil
			},
			Logf: func(string, ...any) {},
		}

		cmd := NewRootCommand(runner)
		cmd.SetArgs([]string{
			"--address", "1.2.3.4",
			"--ports", "80",
			"--timeout", "250ms",
		})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 scan call, got %d", calls)
		}
		if lastHost != "[1.2.3.4]" {
			t.Fatalf("expected host [1.2.3.4], got %q", lastHost)
		}
		if lastPort != 80 {
			t.Fatalf("expected port 80, got %d", lastPort)
		}
		if lastTimeout != 250*time.Millisecond {
			t.Fatalf("expected timeout 250ms, got %v", lastTimeout)
		}
	})

	t.Run("missing ports", func(t *testing.T) {
		cmd := NewRootCommand(app.Runner{})
		cmd.SetArgs([]string{"--address", "1.2.3.4"})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		err := cmd.Execute()
		if !errors.Is(err, app.ErrNoPorts) {
			t.Fatalf("expected ErrNoPorts, got %v", err)
		}
	})

	t.Run("missing address and domain", func(t *testing.T) {
		cmd := NewRootCommand(app.Runner{})
		cmd.SetArgs([]string{"--ports", "80"})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		err := cmd.Execute()
		if !errors.Is(err, app.ErrNoAddressOrDomain) {
			t.Fatalf("expected ErrNoAddressOrDomain, got %v", err)
		}
	})

	t.Run("invalid ports", func(t *testing.T) {
		cmd := NewRootCommand(app.Runner{})
		cmd.SetArgs([]string{"--address", "1.2.3.4", "--ports", "bad"})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		err := cmd.Execute()
		if err == nil || !strings.Contains(err.Error(), "invalid port number") {
			t.Fatalf("expected invalid port error, got %v", err)
		}
	})
}
