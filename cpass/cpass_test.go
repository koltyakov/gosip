package cpass

import (
	"testing"
)

func TestUsingMachineID(t *testing.T) {
	const secret = "secret"
	c := Cpass("")

	encoded, err := c.Encode(secret)
	if err != nil {
		t.Error(err)
	}
	if secret == encoded {
		t.Error("got encryption error")
	}

	decoded, err := c.Decode(encoded)
	if err != nil {
		t.Error(err)
	}
	if secret != decoded {
		t.Error("got decription error")
	}
}

func TestCustomEncryptionKey(t *testing.T) {
	const secret = "secret"
	c1 := Cpass("")
	c2 := Cpass("CUSTOM_KEY")

	encoded1, err := c1.Encode(secret)
	if err != nil {
		t.Error(err)
	}

	encoded2, err := c2.Encode(secret)
	if err != nil {
		t.Error(err)
	}

	if encoded1 == encoded2 {
		t.Error("got encryption error")
	}
}
