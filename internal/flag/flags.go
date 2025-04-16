package flag

import (
	"flag"
	"fmt"
	"time"
)

type Flags struct {
	Address      string
	Domain       string
	Ports        []uint16
	ForceUseIPv4 bool
	ForceUseIPv6 bool
	Timeout      time.Duration
	IgnoreErrors bool
}

func Parse() (Flags, error) {
	var (
		address = flag.String("address", "", "Target IP address")
		ports   sliceInt
		domain  = flag.String("domain", "", "Target domain name")
		useIPv4 = flag.Bool("4", false, "Force use IPv4 protocol")
		useIPv6 = flag.Bool("6", false, "Force use IPv6 protocol")
		timeout = flag.Duration("timeout", 10*time.Millisecond, "Maximum time to establish a connection")
		ignore  = flag.Bool("ignore-errors", false, "Ignore errors when establishing a connection")
	)

	flag.Var(&ports, "ports", "Comma-separated list of ports (0-65535 range)")
	flag.Parse()

	if *address == "" && *domain == "" {
		return Flags{}, fmt.Errorf("no address or domain specified")
	}

	if len(ports) == 0 {
		return Flags{}, fmt.Errorf("no ports specified")
	}

	return Flags{
		Address:      *address,
		Domain:       *domain,
		Ports:        ports,
		ForceUseIPv4: *useIPv4,
		ForceUseIPv6: *useIPv6,
		Timeout:      *timeout,
		IgnoreErrors: *ignore,
	}, nil
}
