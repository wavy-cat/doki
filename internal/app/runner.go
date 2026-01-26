package app

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/wavy-cat/doki/pkg/lookup"
	"github.com/wavy-cat/doki/pkg/scanner"
)

var (
	ErrNoAddressOrDomain = errors.New("no address or domain specified")
	ErrNoPorts           = errors.New("no ports specified")
)

type Options struct {
	Address      string
	Domain       string
	Ports        []uint16
	ForceUseIPv4 bool
	ForceUseIPv6 bool
	Timeout      time.Duration
	IgnoreErrors bool
}

type Runner struct {
	LookupA    func(string) (string, error)
	LookupAAAA func(string) (string, error)
	LookupIP   func(string) (string, error)
	ScanPort   func(string, uint16, time.Duration) error
	Logf       func(string, ...any)
}

func DefaultRunner() Runner {
	return Runner{
		LookupA:    lookup.GetARecord,
		LookupAAAA: lookup.GetAAAARecord,
		LookupIP:   lookup.GetIPRecord,
		ScanPort:   scanner.ScanPort,
		Logf:       log.Printf,
	}
}

func (r Runner) withDefaults() Runner {
	if r.LookupA == nil || r.LookupAAAA == nil || r.LookupIP == nil || r.ScanPort == nil || r.Logf == nil {
		defaults := DefaultRunner()
		if r.LookupA == nil {
			r.LookupA = defaults.LookupA
		}
		if r.LookupAAAA == nil {
			r.LookupAAAA = defaults.LookupAAAA
		}
		if r.LookupIP == nil {
			r.LookupIP = defaults.LookupIP
		}
		if r.ScanPort == nil {
			r.ScanPort = defaults.ScanPort
		}
		if r.Logf == nil {
			r.Logf = defaults.Logf
		}
	}
	return r
}

func (r Runner) Run(opts Options) error {
	if opts.Address == "" && opts.Domain == "" {
		return ErrNoAddressOrDomain
	}
	if len(opts.Ports) == 0 {
		return ErrNoPorts
	}

	r = r.withDefaults()

	address := opts.Address
	if address == "" {
		var err error
		switch {
		case opts.ForceUseIPv4:
			address, err = r.LookupA(opts.Domain)
		case opts.ForceUseIPv6:
			address, err = r.LookupAAAA(opts.Domain)
		default:
			address, err = r.LookupIP(opts.Domain)
		}
		if err != nil {
			return fmt.Errorf("error when retrieving a record from DNS: %w", err)
		}
	}

	address = fmt.Sprintf("[%s]", address)
	for _, port := range opts.Ports {
		err := r.ScanPort(address, port, opts.Timeout)
		if err != nil && !opts.IgnoreErrors {
			r.Logf("error during connection establishment: %v", err)
		}
	}

	return nil
}
