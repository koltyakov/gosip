package cpass

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

var anchor = "cpass|"

// Encrypt encodes a string value to a locally decodable hash.
func encrypt(decoded string, key []byte) (string, error) {
	plainText := []byte(anchor + decoded)
	block, err := aes.NewCipher(key)
	if err != nil {
		return decoded, err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return decoded, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	encoded := base64.URLEncoding.EncodeToString(cipherText)
	return encoded, nil
}

// Decrypt decodes a locally decodable hash to the original string.
func decrypt(encoded string, key []byte) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return encoded, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return encoded, err
	}
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return encoded, err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	decoded := string(cipherText)

	// By design decrypt with incorrect key must end up with the value
	if strings.Index(decoded, anchor) != 0 {
		return encoded, nil
	}

	decoded = strings.Replace(decoded, anchor, "", 1) // remove anchor from string
	return decoded, nil
}
