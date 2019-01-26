package cpass

import (
	"encoding/hex"
	"fmt"
)

// Crypter - cpass module structure
type Crypter struct {
	encryptionKey []byte
}

// Cpass constructor function.
func Cpass(masterKey string) *Crypter {
	if masterKey == "" {
		key, err := getMachineID(false)
		if err != nil {
			masterKey = "CPASS_EMPTY_KEY" // TODO: Fallback logic
		} else {
			masterKey = key
		}
	}
	fmt.Println(hex.EncodeToString(hashCipherKey(masterKey)))
	return &Crypter{hashCipherKey(masterKey)}
}

// Encode encodes a string value to a locally decodable hash.
func (c *Crypter) Encode(data string) (string, error) {
	return encrypt(data, c.encryptionKey)
}

// Decode decodes a locally decodable hash to the original string.
func (c *Crypter) Decode(data string) (string, error) {
	return decrypt(data, c.encryptionKey)
}
