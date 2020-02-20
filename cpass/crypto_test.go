package cpass

import (
	"strings"
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

func TestEdgeCases(t *testing.T) {

	t.Run("encrypt/EmptyKey", func(t *testing.T) {
		if _, err := encrypt("secret", []byte("")); err != nil {
			if !strings.Contains(err.Error(), "invalid key size") {
				t.Error("empty key should not pass")
			}
		}
	})

	t.Run("encrypt/WrongSizeKey", func(t *testing.T) {
		if _, err := encrypt("secret", []byte("wrong_size")); err != nil {
			if !strings.Contains(err.Error(), "invalid key size") {
				t.Error("short key should not pass")
			}
		}
	})

	t.Run("decrypt/EmptyKey", func(t *testing.T) {
		secured, err := encrypt("secret", hashCipherKey("key"))
		if err != nil {
			t.Error(err)
		}
		if _, err := decrypt(secured, []byte("")); err != nil {
			if !strings.Contains(err.Error(), "invalid key size") {
				t.Error("empty key should not pass")
			}
		}
	})

	t.Run("decrypt/IllegaBase64", func(t *testing.T) {
		if _, err := decrypt("incorrect", []byte("key")); err != nil {
			if !strings.Contains(err.Error(), "illegal base64 data") {
				t.Error("illegal base64 data should not pass")
			}
		}
	})

	t.Run("decrypt/SmallBlockSize", func(t *testing.T) {
		if _, err := decrypt("YQ==", hashCipherKey("key")); err != nil {
			if !strings.Contains(err.Error(), "ciphertext block size is too short") {
				t.Error("too short ciphertext block size should not pass")
			}
		}
	})

}
