package cpass

import (
	"crypto/hmac"
	"crypto/sha256"
)

// CipherKey must be 32 chars long because block size is 16 bytes
func hashCipherKey(key string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func getHmac(text, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(text)
	return mac.Sum(nil)
}
