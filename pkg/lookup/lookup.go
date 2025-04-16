package lookup

import (
	"errors"
	"fmt"
	"net"
)

var ErrNoRecords = fmt.Errorf("there are no records on the domain")

var lookupIP = net.LookupIP

// GetARecord возвращает A (IPv4) запись.
// Если получено несколько записей, то отдаётся первая из них.
func GetARecord(domain string) (string, error) {
	ips, err := lookupIP(domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), err
		}
	}
	return "", ErrNoRecords
}

// GetAAAARecord возвращает AAAA (IPv6) запись.
// Если получено несколько записей, то отдаётся первая из них.
func GetAAAARecord(domain string) (string, error) {
	ips, err := lookupIP(domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip.To4() == nil && ip.To16() != nil {
			return ip.String(), err
		}
	}
	return "", ErrNoRecords
}

// GetIPRecord возвращает AAAA (IPv6), либо A (IPv4) запись.
// Если получено несколько записей, то отдаётся первая из них, преимущественно A.
func GetIPRecord(domain string) (string, error) {
	ip, err := GetARecord(domain)
	if err == nil {
		return ip, nil
	}
	if !errors.Is(err, ErrNoRecords) {
		return "", err
	}

	ip, err = GetAAAARecord(domain)
	if err == nil {
		return ip, nil
	}
	if !errors.Is(err, ErrNoRecords) {
		return "", err
	}

	return "", ErrNoRecords
}
