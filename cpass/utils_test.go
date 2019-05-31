package cpass

import (
	"testing"
)

func TestCipherHasher(t *testing.T) {
	b1 := hashCipherKey("")
	b2 := hashCipherKey("any value")

	if len(b1) != 32 {
		t.Error("got incorrect value, must be 32 chars long")
	}

	if len(b2) != 32 {
		t.Error("got incorrect value, must be 32 chars long")
	}
}

func TestHashingDiffValues(t *testing.T) {
	b1 := hashCipherKey("one value")
	b2 := hashCipherKey("another value")

	if string(b1) == string(b2) {
		t.Error("hashing failure, same hashes for different keys")
	}
}


func TestHashingSameValues(t *testing.T) {
	b1 := hashCipherKey("same value")
	b2 := hashCipherKey("same value")

	if string(b1) != string(b2) {
		t.Error("hashing failure, different hashes for the same values")
	}
}
