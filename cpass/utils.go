package cpass

import (
	"crypto/hmac"
	"crypto/sha256"
)

// CipherKey must be 32 chars long because block size is 16 bytes
func hashCipherKey(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return string(hasher.Sum(nil)) // hex.EncodeToString(
}

func getHmac(text, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(text)
	return mac.Sum(nil)
}

func pad(buf []byte, size int) []byte {
	bufLen := len(buf)
	padLen := size - bufLen%size
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(padLen)
	}
	return padded
}

func unpad(padded []byte, size int) []byte {
	if len(padded)%size != 0 {
		panic("Padded value has incorrect size.")
	}
	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf
}
