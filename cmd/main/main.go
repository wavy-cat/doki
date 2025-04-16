package main

import (
	"fmt"
	"github.com/wavy-cat/doki/internal/flag"
	"github.com/wavy-cat/doki/pkg/lookup"
	"github.com/wavy-cat/doki/pkg/scanner"
	"log"
)

func main() {
	flags, err := flag.Parse()
	if err != nil {
		log.Fatal(err)
	}

	address := flags.Address
	if address == "" {
		address = flags.Domain
		switch {
		case flags.ForceUseIPv4:
			address, err = lookup.GetARecord(flags.Domain)
		case flags.ForceUseIPv6:
			address, err = lookup.GetAAAARecord(flags.Domain)
		default:
			address, err = lookup.GetIPRecord(flags.Domain)
		}

		address = fmt.Sprintf("[%s]", address)

		if err != nil {
			log.Fatalf("error when retrieving a record from DNS: %v", err)
		}
	}

	for _, port := range flags.Ports {
		err = scanner.ScanPort(address, port, flags.Timeout)
		if err != nil && !flags.IgnoreErrors {
			log.Printf("error during connection establishment: %v", err)
		}
	}
}
