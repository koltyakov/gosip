package main

import (
	"encoding/base64"
	"fmt"

	"github.com/koltyakov/gosip/aes256cbc"
)

func main() {

	original := "My string"
	fmt.Printf("original: %s\n", original)

	encoded, _ := aes256cbc.Encrypt(original)
	fmt.Printf("encoded: %s\n", encoded)

	decoded, _ := aes256cbc.Decrypt(encoded)
	fmt.Printf("decoded: %s\n", decoded)

	// decodedFromNode, _ := aes256cbc.Decrypt("81aa2f840f0cc3a470cc070956e64924237d4fa9aa05057f7e2cf26cd4c6edd5")
	decodedFromNode, _ := aes256cbc.Decrypt("6e60af6d34cbdfce11a9c8c09cf0f67b5bcb0e3fb16cd569dc6de407d17ebfd3")
	fmt.Printf("decodedFromNode: %s\n", decodedFromNode)

	// h, _ := hex.DecodeString("ab62b3220d3b86bbe4b469f7f2463444880db70d223b9f1f93b7b15c02a5366d")
	h := []byte("ab62b3220d3b86bbe4b469f7f2463444880db70d223b9f1f93b7b15c02a5366d")
	k := base64.StdEncoding.EncodeToString(h)
	fmt.Printf("k: %s\n", k)

	b, _ := base64.StdEncoding.DecodeString("NmU2MGFmNmQzNGNiZGZjZTExYTljOGMwOWNmMGY2N2I1YmNiMGUzZmIxNmNkNTY5ZGM2ZGU0MDdkMTdlYmZkMw==")
	fmt.Printf("b: %s\n", b)

}

// go run cmd/test/main.go
