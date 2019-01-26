package cpass

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encodes a string value to a locally decodable hash.
func encrypt(decoded string, key []byte) (encoded string, err error) {
	plainText := []byte(decoded)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	encoded = base64.URLEncoding.EncodeToString(cipherText)
	return
}

// Decrypt decodes a locally decodable hash to the original string.
func decrypt(encoded string, key []byte) (decoded string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short")
		return
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	decoded = string(cipherText)
	return
}
