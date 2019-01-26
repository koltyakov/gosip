package cpass

import "testing"

func TestUsingMachineID(t *testing.T) {
	const secret = "secret"
	c := Cpass("")

	encoded, err := c.Encode(secret)
	if err != nil {
		t.Error(err)
	}
	if secret == encoded {
		t.Error("Got encryption error")
	}

	decoded, err := c.Decode(encoded)
	if err != nil {
		t.Error(err)
	}
	if secret != decoded {
		t.Error("Got decription error")
	}
}
