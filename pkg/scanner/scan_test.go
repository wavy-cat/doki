package scanner_test

import (
	"github.com/wavy-cat/doki/pkg/scanner"
	"testing"
	"time"
)

func TestScanner(t *testing.T) {
	t.Run("Knock 1.1.1.1", func(t *testing.T) {
		t.Parallel()

		err := scanner.ScanPort("1.1.1.1", 443, time.Second)
		if err != nil {
			t.Fatalf("ScanPort returned error: %v", err)
		}
	})

	t.Run("Knock error", func(t *testing.T) {
		t.Parallel()

		err := scanner.ScanPort("256.0.0.0", 0, time.Second)
		if err == nil {
			t.Fatalf("ScanPort = nil, want error")
		}
	})
}
