package cmd

import (
	"testing"
)

func TestPortsFlagSet(t *testing.T) {
	t.Run("single value", func(t *testing.T) {
		var p portsFlag
		if err := p.Set("80"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(p) != 1 || p[0] != 80 {
			t.Fatalf("unexpected ports: %v", p)
		}
	})

	t.Run("multiple values with spaces", func(t *testing.T) {
		var p portsFlag
		if err := p.Set("80, 443,3000"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []uint16{80, 443, 3000}
		for i, v := range expected {
			if p[i] != v {
				t.Fatalf("expected %v, got %v", expected, p)
			}
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		var p portsFlag
		if err := p.Set("bad"); err == nil {
			t.Fatal("expected error for invalid port")
		}
	})

	t.Run("empty value", func(t *testing.T) {
		var p portsFlag
		if err := p.Set(""); err == nil {
			t.Fatal("expected error for empty port value")
		}
	})

	t.Run("empty entry", func(t *testing.T) {
		var p portsFlag
		if err := p.Set("80,,443"); err == nil {
			t.Fatal("expected error for empty port entry")
		}
	})
}

func TestPortsFlagString(t *testing.T) {
	p := portsFlag{80, 443}
	if p.String() != "80,443" {
		t.Fatalf("unexpected string: %q", p.String())
	}

	var empty portsFlag
	if empty.String() != "" {
		t.Fatalf("expected empty string, got %q", empty.String())
	}

	var nilPorts *portsFlag
	if nilPorts.String() != "" {
		t.Fatalf("expected empty string for nil ports, got %q", nilPorts.String())
	}
}

func TestPortsFlagType(t *testing.T) {
	var p portsFlag
	if p.Type() != "ports" {
		t.Fatalf("unexpected type: %q", p.Type())
	}
}
