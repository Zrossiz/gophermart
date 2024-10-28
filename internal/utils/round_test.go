package utils

import "testing"

func TestRound(t *testing.T) {
	rounded := Round(123.1123123, 2)
	if rounded != 123.11 {
		t.Fatal("expected 123.11, got: ", rounded)
	}
}
