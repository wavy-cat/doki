package flag

import (
	"reflect"
	"testing"
)

func TestSliceInt_Set(t *testing.T) {
	t.Run("single number", func(t *testing.T) {
		t.Parallel()
		var s sliceInt
		err := s.Set("8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := sliceInt{8080}
		if !reflect.DeepEqual(s, expected) {
			t.Fatalf("expected %v, got %v", expected, s)
		}
	})

	t.Run("multiple numbers", func(t *testing.T) {
		t.Parallel()
		var s sliceInt
		err := s.Set("80,443,3000")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := sliceInt{80, 443, 3000}
		if !reflect.DeepEqual(s, expected) {
			t.Fatalf("expected %v, got %v", expected, s)
		}
	})

	t.Run("invalid number", func(t *testing.T) {
		t.Parallel()
		var s sliceInt
		err := s.Set("notanumber")
		if err == nil {
			t.Fatal("expected error for invalid number")
		}
	})

	t.Run("mixed valid and invalid", func(t *testing.T) {
		t.Parallel()
		var s sliceInt
		err := s.Set("1234,hello,5678")
		if err == nil {
			t.Fatal("expected error for mixed values")
		}
	})

	t.Run("get string", func(t *testing.T) {
		t.Parallel()
		s := sliceInt{80, 443}
		expected := "[80 443]"
		if s.String() != expected {
			t.Fatalf("expected string %q, got %q", expected, s.String())
		}
	})
}
