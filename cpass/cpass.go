package cpass

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// Cpass a simple encryption library.
type Cpass struct {
	encryptionKey []byte
}

// NewCpass is Cpass constructor function.
func NewCpass(masterKey string) *Cpass {
	if masterKey == "" {
		masterKey = hashCipherKey(getMachineID(false))
	}
	return &Cpass{[]byte(masterKey)}
}

// Encode encodes a string value to a locally decodable hash.
func (c *Cpass) Encode(data string) (string, error) {
	step1, err := c.Encrypt([]byte(data))
	if err != nil {
		return data, err
	}
	step2 := hex.EncodeToString(step1)
	step3 := []byte(step2)
	step4 := base64.StdEncoding.EncodeToString(step3)
	return step4, nil
}

// Decode decodes a locally decodable hash to the original string.
func (c *Cpass) Decode(data string) (string, error) {
	step1, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data, err
	}
	step2 := fmt.Sprintf("%s", step1)
	step3, err := hex.DecodeString(step2)
	if err != nil {
		return data, err
	}
	step4, err := c.Decrypt([]byte(step3))
	if err != nil {
		return data, err
	}
	step5 := string(step4)
	return step5, nil
}

// Encrypt encodes a string value to a locally decodable hash.
func (c *Cpass) Encrypt(data []byte) ([]byte, error) {
	data = pad(data, aes.BlockSize)
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return []byte{}, err
	}
	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], data)
	return cipherText, nil
}

// Decrypt decodes a locally decodable hash to the original string.
func (c *Cpass) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		panic(err)
	}
	if len(data) < aes.BlockSize {
		panic("data too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	if len(data)%aes.BlockSize != 0 {
		panic("data is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	data = unpad(data, aes.BlockSize)
	return data, nil
}
