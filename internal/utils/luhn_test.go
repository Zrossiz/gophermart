package utils

import "testing"

func TestLuhn(t *testing.T) {
	validLuhnNum := "12345678903"
	invalidLuhnNum := "1234567890"

	valid := IsLuhn(validLuhnNum)
	if !valid {
		t.Fatal("expected true, got false")
	}

	invalid := IsLuhn(invalidLuhnNum)
	if invalid {
		t.Fatal("expected false, got true")
	}
}
