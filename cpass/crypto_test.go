package cpass

import (
	"testing"
)

var masterKey = hashCipherKey("MY_MASTER_KEY")
var incorrectKey = hashCipherKey("INCORRECT_MASTER_KEY")

func TestEncryptAndDecrypt(t *testing.T) {
	rawSecret := "secret"

	encrypted, err := encrypt(rawSecret, masterKey)
	if err != nil {
		t.Error(err)
	}

	if rawSecret == encrypted {
		t.Error("encrypted is equal to raw secret")
	}

	if len(encrypted) <= len(rawSecret) {
		t.Error("encrypted is shorter than raw secret")
	}

	decrypted, err := decrypt(encrypted, masterKey)
	if err != nil {
		t.Error(err)
	}

	if decrypted == encrypted {
		t.Error("encrypted should not be equal to raw secret")
	}

	if decrypted != rawSecret {
		t.Error("decrypted should not be equal to raw secret")
	}
}

func TestMultiEncryption(t *testing.T) {
	rawSecret := "secret"

	encrypted1, err := encrypt(rawSecret, masterKey)
	if err != nil {
		t.Error(err)
	}

	encrypted2, err := encrypt(rawSecret, masterKey)
	if err != nil {
		t.Error(err)
	}

	if encrypted1 == encrypted2 {
		t.Error("encrypted values should not be the same")
	}
}

func TestDecryptWithIncorrectKey(t *testing.T) {
	rawSecret := "secret"

	encrypted, err := encrypt(rawSecret, masterKey)
	if err != nil {
		t.Error(err)
	}

	if encrypted == rawSecret {
		t.Error("encryption failed")
	}

	decrypted, err := decrypt(encrypted, incorrectKey)
	if err != nil {
		t.Error(err)
	}

	if encrypted != decrypted {
		t.Error("decrypted by incorrect key is not equal to secret")
	}
}
