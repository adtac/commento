package main

import (
	"testing"
)

func TestRandomHexBasics(t *testing.T) {
	hex1, err := randomHex(32)
	if err != nil {
		t.Errorf("unexpected error creating hex: %v", err)
		return
	}

	if hex1 == "" {
		t.Errorf("randomly generated hex empty")
		return
	}

	hex2, err := randomHex(32)
	if err != nil {
		t.Errorf("unexpected error creating hex: %v", err)
		return
	}

	if hex1 == hex2 {
		t.Errorf("two randomly generated hexes found to be the same: '%s'", hex1)
		return
	}
}
