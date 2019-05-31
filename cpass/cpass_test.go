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
		t.Error("got decryption error")
	}
}

func TestCustomEncryptionKey(t *testing.T) {
	const secret = "secret"
	c1 := Cpass("")
	c2 := Cpass("CUSTOM_KEY")

	encoded, err := c1.Encode(secret)
	if err != nil {
		t.Error(err)
	}

	decoded, err := c2.Decode(encoded)
	if err != nil {
		t.Error(err)
	}

	if secret == decoded {
		t.Error("got encryption error")
	}
}

func TestEmptyShouldUseMachineID(t *testing.T) {
	const secret = "secret"

	machineID, err := getMachineID(false)
	if err != nil {
		t.Error(err)
	}

	c1 := Cpass(machineID)
	c2 := Cpass("")

	encoded, err := c1.Encode(secret)
	if err != nil {
		t.Error(err)
	}

	decoded, err := c2.Decode(encoded)
	if err != nil {
		t.Error(err)
	}

	if secret != decoded {
		t.Error("got encryption error")
	}
}

func TestMasterKeyIsNotEmpty(t *testing.T) {
	c := Cpass("")
	if string(c.encryptionKey) == string(hashCipherKey("")) {
		t.Error("got master key helper error")
	}
}
