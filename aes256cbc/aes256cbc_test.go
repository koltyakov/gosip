package aes256cbc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	t.Run("Encrypts and decrypts", func(t *testing.T) {
		plainTexts := []string{"1234567890", "123456789012345678901234567890123456789012345678901234567890", "1", ""}
		for _, plainText := range plainTexts {
			encrypted, err := Encrypt(plainText)
			if err != nil {
				t.Fatalf("Failed to encrypt: %s - %s", plainText, err.Error())
			}
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Failed to decrypt: %s - %s", plainText, err.Error())
			}
			assert.Equal(t, plainText, decrypted)
		}
	})
}
