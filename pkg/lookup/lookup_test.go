package lookup

import (
	"errors"
	"net"
	"testing"
)

func setMockLookupIP(f func(string) ([]net.IP, error)) func() {
	original := lookupIP
	lookupIP = f
	return func() { lookupIP = original }
}

func TestLookup(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(string) ([]net.IP, error)
		testFunc func(t *testing.T)
	}{
		{
			name: "A record exists",
			mock: func(string) ([]net.IP, error) {
				return []net.IP{net.ParseIP("192.0.2.1")}, nil
			},
			testFunc: func(t *testing.T) {
				ip, err := GetARecord("example.com")
				if err != nil || ip != "192.0.2.1" {
					t.Fatalf("expected 192.0.2.1, got %v, err: %v", ip, err)
				}
			},
		},
		{
			name: "AAAA record exists",
			mock: func(string) ([]net.IP, error) {
				return []net.IP{net.ParseIP("2001:db8::1")}, nil
			},
			testFunc: func(t *testing.T) {
				ip, err := GetAAAARecord("example.com")
				if err != nil || ip != "2001:db8::1" {
					t.Fatalf("expected 2001:db8::1, got %v, err: %v", ip, err)
				}
			},
		},
		{
			name: "A record found",
			mock: func(string) ([]net.IP, error) {
				return []net.IP{net.ParseIP("192.0.2.1")}, nil
			},
			testFunc: func(t *testing.T) {
				ip, err := GetIPRecord("example.com")
				if err != nil || ip != "192.0.2.1" {
					t.Fatalf("expected 192.0.2.1, got %v, err: %v", ip, err)
				}
			},
		},
		{
			name: "A record not found, fallback to AAAA",
			mock: func(string) ([]net.IP, error) {
				return []net.IP{net.ParseIP("2001:db8::1")}, nil
			},
			testFunc: func(t *testing.T) {
				ip, err := GetIPRecord("example.com")
				if err != nil || ip != "2001:db8::1" {
					t.Fatalf("expected 2001:db8::1 fallback, got %v, err: %v", ip, err)
				}
			},
		},
		{
			name: "Only IPv6 returned for A, should fail",
			mock: func(string) ([]net.IP, error) {
				return []net.IP{net.ParseIP("2001:db8::1")}, nil
			},
			testFunc: func(t *testing.T) {
				_, err := GetARecord("example.com")
				if !errors.Is(err, ErrNoRecords) {
					t.Fatalf("expected ErrNoRecords, got %v", err)
				}
			},
		},
		{
			name: "Only IPv4 returned for AAAA, should fail",
			mock: func(string) ([]net.IP, error) {
				return []net.IP{net.ParseIP("192.0.2.1")}, nil
			},
			testFunc: func(t *testing.T) {
				_, err := GetAAAARecord("example.com")
				if !errors.Is(err, ErrNoRecords) {
					t.Fatalf("expected ErrNoRecords, got %v", err)
				}
			},
		},
		{
			name: "Records not found, should fail",
			mock: func(string) ([]net.IP, error) {
				return nil, nil
			},
			testFunc: func(t *testing.T) {
				_, err := GetIPRecord("example.com")
				if !errors.Is(err, ErrNoRecords) {
					t.Fatalf("expected ErrNoRecords, got %v", err)
				}
			},
		},
		{
			name: "Lookup fails with error (A)",
			mock: func(string) ([]net.IP, error) {
				return nil, errors.New("dns failure")
			},
			testFunc: func(t *testing.T) {
				_, err := GetARecord("example.com")
				if err == nil || err.Error() != "dns failure" {
					t.Fatalf("expected dns failure, got %v", err)
				}
			},
		},
		{
			name: "Lookup fails with error (AAAA)",
			mock: func(string) ([]net.IP, error) {
				return nil, errors.New("dns failure")
			},
			testFunc: func(t *testing.T) {
				_, err := GetAAAARecord("example.com")
				if err == nil || err.Error() != "dns failure" {
					t.Fatalf("expected dns failure, got %v", err)
				}
			},
		},
		{
			name: "Lookup fails with error (IP)",
			mock: func(string) ([]net.IP, error) {
				return nil, errors.New("dns failure")
			},
			testFunc: func(t *testing.T) {
				_, err := GetIPRecord("example.com")
				if err == nil || err.Error() != "dns failure" {
					t.Fatalf("expected dns failure, got %v", err)
				}
			},
		},
		{
			name: "A record not found, AAAA lookup fails",
			// Здесь мы создаём замыкание с переменной-счётчиком,
			// которая сохраняется между вызовами lookupIP.
			mock: func() func(string) ([]net.IP, error) {
				var callCount int
				return func(domain string) ([]net.IP, error) {
					if callCount == 0 {
						callCount++
						// При первом вызове (для GetARecord)
						// возвращаем IPv6, чтобы A-запись не нашлась.
						return []net.IP{net.ParseIP("2001:db8::1")}, nil
					}
					// При втором вызове (для GetAAAARecord) имитируем ошибку.
					return nil, errors.New("no AAAA")
				}
			}(),
			testFunc: func(t *testing.T) {
				_, err := GetIPRecord("example.com")
				if err == nil || err.Error() != "no AAAA" {
					t.Fatalf("expected error 'no AAAA', got %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc // захват в локальную переменную
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			restore := setMockLookupIP(tc.mock)
			defer restore()
			tc.testFunc(t)
		})
	}
}
