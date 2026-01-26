package app

import (
	"errors"
	"testing"
	"time"
)

func TestRunnerRunValidation(t *testing.T) {
	runner := Runner{}

	err := runner.Run(Options{Ports: []uint16{80}})
	if !errors.Is(err, ErrNoAddressOrDomain) {
		t.Fatalf("expected ErrNoAddressOrDomain, got %v", err)
	}

	err = runner.Run(Options{Address: "1.2.3.4"})
	if !errors.Is(err, ErrNoPorts) {
		t.Fatalf("expected ErrNoPorts, got %v", err)
	}
}

func TestRunnerWithDefaults(t *testing.T) {
	var runner Runner
	runner = runner.withDefaults()

	if runner.LookupA == nil {
		t.Fatal("expected LookupA to be set")
	}
	if runner.LookupAAAA == nil {
		t.Fatal("expected LookupAAAA to be set")
	}
	if runner.LookupIP == nil {
		t.Fatal("expected LookupIP to be set")
	}
	if runner.ScanPort == nil {
		t.Fatal("expected ScanPort to be set")
	}
	if runner.Logf == nil {
		t.Fatal("expected Logf to be set")
	}
}

func TestRunnerRunAddressProvided(t *testing.T) {
	var (
		calls       int
		lastHost    string
		lastPort    uint16
		lastTimeout time.Duration
	)

	runner := Runner{
		LookupA: func(string) (string, error) {
			t.Fatal("LookupA should not be called when address is provided")
			return "", nil
		},
		LookupAAAA: func(string) (string, error) {
			t.Fatal("LookupAAAA should not be called when address is provided")
			return "", nil
		},
		LookupIP: func(string) (string, error) {
			t.Fatal("LookupIP should not be called when address is provided")
			return "", nil
		},
		ScanPort: func(host string, port uint16, timeout time.Duration) error {
			calls++
			lastHost = host
			lastPort = port
			lastTimeout = timeout
			return nil
		},
		Logf: func(string, ...any) {},
	}

	err := runner.Run(Options{
		Address: "1.2.3.4",
		Ports:   []uint16{80},
		Timeout: 150 * time.Millisecond,
	})
	if err != nil {
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
	if lastTimeout != 150*time.Millisecond {
		t.Fatalf("expected timeout 150ms, got %v", lastTimeout)
	}
}

func TestRunnerRunDomainLookups(t *testing.T) {
	t.Run("force IPv4", func(t *testing.T) {
		lookupCalled := 0
		runner := Runner{
			LookupA: func(domain string) (string, error) {
				lookupCalled++
				if domain != "example.com" {
					t.Fatalf("unexpected domain: %q", domain)
				}
				return "192.0.2.1", nil
			},
			LookupAAAA: func(string) (string, error) {
				t.Fatal("LookupAAAA should not be called")
				return "", nil
			},
			LookupIP: func(string) (string, error) {
				t.Fatal("LookupIP should not be called")
				return "", nil
			},
			ScanPort: func(host string, port uint16, timeout time.Duration) error {
				if host != "[192.0.2.1]" {
					t.Fatalf("unexpected host: %q", host)
				}
				return nil
			},
			Logf: func(string, ...any) {},
		}

		err := runner.Run(Options{
			Domain:       "example.com",
			Ports:        []uint16{443},
			ForceUseIPv4: true,
			Timeout:      time.Second,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lookupCalled != 1 {
			t.Fatalf("expected LookupA to be called once, got %d", lookupCalled)
		}
	})

	t.Run("force IPv6", func(t *testing.T) {
		lookupCalled := 0
		runner := Runner{
			LookupA: func(string) (string, error) {
				t.Fatal("LookupA should not be called")
				return "", nil
			},
			LookupAAAA: func(domain string) (string, error) {
				lookupCalled++
				if domain != "example.com" {
					t.Fatalf("unexpected domain: %q", domain)
				}
				return "2001:db8::1", nil
			},
			LookupIP: func(string) (string, error) {
				t.Fatal("LookupIP should not be called")
				return "", nil
			},
			ScanPort: func(host string, port uint16, timeout time.Duration) error {
				if host != "[2001:db8::1]" {
					t.Fatalf("unexpected host: %q", host)
				}
				return nil
			},
			Logf: func(string, ...any) {},
		}

		err := runner.Run(Options{
			Domain:       "example.com",
			Ports:        []uint16{443},
			ForceUseIPv6: true,
			Timeout:      time.Second,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lookupCalled != 1 {
			t.Fatalf("expected LookupAAAA to be called once, got %d", lookupCalled)
		}
	})

	t.Run("default lookup", func(t *testing.T) {
		lookupCalled := 0
		runner := Runner{
			LookupA: func(string) (string, error) {
				t.Fatal("LookupA should not be called")
				return "", nil
			},
			LookupAAAA: func(string) (string, error) {
				t.Fatal("LookupAAAA should not be called")
				return "", nil
			},
			LookupIP: func(domain string) (string, error) {
				lookupCalled++
				if domain != "example.com" {
					t.Fatalf("unexpected domain: %q", domain)
				}
				return "192.0.2.1", nil
			},
			ScanPort: func(host string, port uint16, timeout time.Duration) error {
				if host != "[192.0.2.1]" {
					t.Fatalf("unexpected host: %q", host)
				}
				return nil
			},
			Logf: func(string, ...any) {},
		}

		err := runner.Run(Options{
			Domain:  "example.com",
			Ports:   []uint16{443},
			Timeout: time.Second,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lookupCalled != 1 {
			t.Fatalf("expected LookupIP to be called once, got %d", lookupCalled)
		}
	})
}

func TestRunnerRunLookupError(t *testing.T) {
	lookupErr := errors.New("lookup failed")
	runner := Runner{
		LookupIP: func(string) (string, error) {
			return "", lookupErr
		},
		ScanPort: func(string, uint16, time.Duration) error {
			t.Fatal("ScanPort should not be called on lookup error")
			return nil
		},
		Logf: func(string, ...any) {},
	}

	err := runner.Run(Options{
		Domain:  "example.com",
		Ports:   []uint16{80},
		Timeout: time.Second,
	})
	if !errors.Is(err, lookupErr) {
		t.Fatalf("expected wrapped lookup error, got %v", err)
	}
}

func TestRunnerRunScanErrors(t *testing.T) {
	t.Run("logs when not ignoring errors", func(t *testing.T) {
		logged := 0
		runner := Runner{
			LookupIP: func(string) (string, error) { return "192.0.2.1", nil },
			ScanPort: func(string, uint16, time.Duration) error {
				return errors.New("scan failed")
			},
			Logf: func(string, ...any) {
				logged++
			},
		}

		err := runner.Run(Options{
			Domain:  "example.com",
			Ports:   []uint16{80},
			Timeout: time.Second,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if logged != 1 {
			t.Fatalf("expected 1 log call, got %d", logged)
		}
	})

	t.Run("does not log when ignoring errors", func(t *testing.T) {
		logged := 0
		runner := Runner{
			LookupIP: func(string) (string, error) { return "192.0.2.1", nil },
			ScanPort: func(string, uint16, time.Duration) error {
				return errors.New("scan failed")
			},
			Logf: func(string, ...any) {
				logged++
			},
		}

		err := runner.Run(Options{
			Domain:       "example.com",
			Ports:        []uint16{80},
			Timeout:      time.Second,
			IgnoreErrors: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if logged != 0 {
			t.Fatalf("expected no log calls, got %d", logged)
		}
	})
}
